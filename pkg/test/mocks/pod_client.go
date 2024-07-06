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

import (
	"context"
	"github.com/telemaco019/duplik8s/pkg/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ListPodsResult struct {
	Objs []core.DuplicableObject
	Err  error
}

func NewListPodResults(pods []string, namespace string, err error) ListPodsResult {
	var objs = make([]core.DuplicableObject, 0)
	for _, pod := range pods {
		objs = append(objs, core.DuplicableObject{
			Name:      pod,
			Namespace: namespace,
		})
	}
	return ListPodsResult{
		Objs: objs,
		Err:  err,
	}
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

func (c *PodClient) ListDuplicable(ctx context.Context, resource schema.GroupVersionResource, namespace string) ([]core.DuplicableObject, error) {
	return c.ListPodsResult.Objs, c.ListPodsResult.Err
}

func (c *PodClient) Duplicate(_ core.DuplicableObject, __ core.PodOverrideOptions) error {
	return c.DuplicatePodResult
}
