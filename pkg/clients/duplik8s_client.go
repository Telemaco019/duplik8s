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

package clients

import (
	"context"
	"github.com/telemaco019/duplik8s/pkg/core"
	"github.com/telemaco019/duplik8s/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Duplik8sClient struct {
	client dynamic.Interface
}

func NewDuplik8sClient(opts utils.KubeOptions) (*Duplik8sClient, error) {
	client, err := utils.NewDynamicClient(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &Duplik8sClient{
		client: client,
	}, nil
}

func (c Duplik8sClient) ListDuplicable(
	ctx context.Context,
	resource schema.GroupVersionResource,
	namespace string,
) ([]core.DuplicableObject, error) {
	unstructuredList, err := c.client.Resource(resource).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var objs []core.DuplicableObject
	for _, u := range unstructuredList.Items {
		objs = append(objs, core.NewDuplicable(u))
	}
	return objs, nil
}

func (c Duplik8sClient) ListDuplicated(ctx context.Context) ([]core.DuplicableObject, error) {
	return nil, nil
}
