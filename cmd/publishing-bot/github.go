/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func githubClient(token string) *github.Client {
	// create github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func ReportOnIssue(e error, logs, token, org, repo string, issue int) error {
	ctx := context.Background()
	client := githubClient(token)

	// filter out token, if it happens to be in the log (it shouldn't!)
	// TODO: Consider using log sanitizer from sigs.k8s.io/release-utils
	logs = strings.ReplaceAll(logs, token, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

	// who am I?
	myself, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get own user: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get own user: HTTP code %d", resp.StatusCode)
	}

	// create new newComment
	body := transfromLogToGithubFormat(logs, 50, fmt.Sprintf("/reopen\n\nThe last publishing run failed: %v", e))

	newComment, resp, err := client.Issues.CreateComment(ctx, org, repo, issue, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return fmt.Errorf("failed to comment on issue #%d: %w", issue, err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to comment on issue #%d: HTTP code %d", issue, resp.StatusCode)
	}

	// delete all other comments from this user
	comments, resp, err := client.Issues.ListComments(ctx, org, repo, issue, &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to get github comments of issue #%d: %w", issue, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get github comments of issue #%d: HTTP code %d", issue, resp.StatusCode)
	}
	for _, c := range comments {
		if *c.User.ID != *myself.ID {
			glog.Infof("Skipping comment %d not by me, but %v", *c.ID, c.User.Name)
			continue
		}
		if *c.ID == *newComment.ID {
			continue
		}

		glog.Infof("Deleting comment %d", *c.ID)
		resp, err = client.Issues.DeleteComment(ctx, org, repo, *c.ID)
		if err != nil {
			return fmt.Errorf("failed to delete github comment %d of issue #%d: %w", *c.ID, issue, err)
		}
		if resp.StatusCode >= 300 {
			return fmt.Errorf("failed to delete github comment %d of issue #%d: HTTP code %d", *c.ID, issue, resp.StatusCode)
		}
	}

	return nil
}

func CloseIssue(token, org, repo string, issue int) error {
	ctx := context.Background()
	client := githubClient(token)

	_, resp, err := client.Issues.Edit(ctx, org, repo, issue, &github.IssueRequest{
		State: github.String("closed"),
	})
	if err != nil {
		return fmt.Errorf("failed to close issue #%d: %w", issue, err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to close issue #%d: HTTP code %d", issue, resp.StatusCode)
	}

	return nil
}

func transfromLogToGithubFormat(original string, maxLines int, headings ...string) string {
	logCount := 0
	transformed := newLogBuilderWithMaxBytes(65000, original).
		AddHeading(headings...).
		AddHeading("```").
		Trim("\n").
		Split("\n").
		Reverse().
		Filter(func(line string) bool {
			if logCount < maxLines {
				if strings.HasPrefix(line, "+") {
					logCount++
				}
				return true
			}
			return false
		}).
		Reverse().
		Join("\n").
		AddTailing("```").
		Log()
	return transformed
}
