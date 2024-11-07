#!/usr/bin/env bash

# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

# shellcheck disable=SC2034
repos=(
    api
    apiextensions-apiserver
    apimachinery
    apiserver
    cli-runtime
    client-go
    cloud-provider
    cluster-bootstrap
    code-generator
    component-base
    component-helpers
    controller-manager
    cri-api
    cri-client
    csi-translation-lib
    dynamic-resource-allocation
    externaljwt
    endpointslice
    kms
    kube-aggregator
    kube-controller-manager
    kube-proxy
    kube-scheduler
    kubectl
    kubelet
    legacy-cloud-providers
    metrics
    mount-utils
    pod-security-admission
    sample-apiserver
    sample-cli-plugin
    sample-controller
)
