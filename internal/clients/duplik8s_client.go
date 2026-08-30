/*
 * Copyright 2025 Michele Zanotti <m.zanotti019@gmail.com>
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
	"fmt"

	"github.com/telemaco019/duplik8s/internal/core"
	"github.com/telemaco019/duplik8s/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

type Duplik8sClient struct {
	dynamic   dynamic.Interface
	discovery discovery.DiscoveryInterface
}

func NewDuplik8sClient(opts utils.KubeOptions) (*Duplik8sClient, error) {
	dynamic, err := utils.NewDynamicClient(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}

	discovery, err := utils.NewDiscoveryClient(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}

	return &Duplik8sClient{
		dynamic:   dynamic,
		discovery: discovery,
	}, nil
}

func (c Duplik8sClient) ListDuplicable(
	ctx context.Context,
	resource schema.GroupVersionResource,
	namespace string,
) ([]core.DuplicableObject, error) {
	unstructuredList, err := c.dynamic.Resource(resource).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var objs []core.DuplicableObject
	for _, u := range unstructuredList.Items {
		objs = append(objs, core.NewDuplicable(u))
	}
	return objs, nil
}

func (c Duplik8sClient) Delete(
	ctx context.Context,
	obj core.DuplicatedObject,
) error {
	// Get a RESTMapper to convert the object's GroupVersionKind into a
	// GroupVersionResource for the dynamic client.
	resources, err := restmapper.GetAPIGroupResources(c.discovery)
	if err != nil {
		return fmt.Errorf("failed to discover API group resources: %w", err)
	}
	restMapper := restmapper.NewDiscoveryRESTMapper(resources)

	mapping, err := restMapper.RESTMapping(
		obj.ObjectKind.GroupVersionKind().GroupKind(),
		obj.ObjectKind.GroupVersionKind().Version,
	)
	if err != nil {
		return fmt.Errorf("failed to map %s to a resource: %w", obj.ObjectKind.GroupVersionKind().String(), err)
	}

	return c.dynamic.Resource(mapping.Resource).Namespace(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
}

func (c Duplik8sClient) ListDuplicated(
	ctx context.Context,
	namespace string,
) ([]core.DuplicatedObject, error) {
	resources := make([]core.DuplicatedObject, 0)
	for _, r := range core.SupportedResources {
		unstructuredList, err := c.dynamic.Resource(r.GVR).
			Namespace(namespace).
			List(ctx, metav1.ListOptions{
				LabelSelector: core.LABEL_DUPLICATED + "=true",
			})
		if err != nil {
			return nil, err
		}

		for _, u := range unstructuredList.Items {
			resources = append(resources, core.DuplicatedObject{
				Name:              u.GetName(),
				Namespace:         u.GetNamespace(),
				ObjectKind:        u.GetObjectKind(),
				CreationTimestamp: u.GetCreationTimestamp(),
			})
		}
	}
	return resources, nil
}
