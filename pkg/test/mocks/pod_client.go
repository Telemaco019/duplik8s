/*
 * Copyright 2024 Michele Zanotti <m.zanotti019@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mocks

import "github.com/telemaco019/duplik8s/pkg/pods"

type ListPodsResult struct {
	Pods []string
	Err  error
}

type PodClient struct {
	ListPodsResult     ListPodsResult
	DuplicatePodResult error
}

func NewPodClient(
	ListPodsResult ListPodsResult,
	DuplicatePodResult error,
) *PodClient {
	return &PodClient{
		ListPodsResult:     ListPodsResult,
		DuplicatePodResult: DuplicatePodResult,
	}
}

func (c *PodClient) ListPods(namespace string) ([]string, error) {
	return c.ListPodsResult.Pods, c.ListPodsResult.Err
}

func (c *PodClient) DuplicatePod(podName string, namespace string, opts pods.PodOverrideOptions) error {
	return c.DuplicatePodResult
}
