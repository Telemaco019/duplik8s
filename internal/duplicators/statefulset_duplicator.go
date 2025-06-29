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

package duplicators

import (
	"context"
	"fmt"
	"github.com/telemaco019/duplik8s/internal/clients"
	"github.com/telemaco019/duplik8s/internal/core"
	"github.com/telemaco019/duplik8s/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StatefulSetClient struct {
	clientset *kubernetes.Clientset
	ctx       context.Context
}

func NewStatefulSetClient(opts utils.KubeOptions) (*StatefulSetClient, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &StatefulSetClient{
		clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

func (c *StatefulSetClient) Duplicate(obj core.DuplicableObject, opts core.DuplicateOpts) error {
	fmt.Printf("duplicating statefulset %s\n", obj.Name)

	// fetch the StatefulSet
	statefulSet, err := c.clientset.AppsV1().StatefulSets(obj.Namespace).Get(c.ctx, obj.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if statefulSet.Labels[core.LABEL_DUPLICATED] == "true" {
		return fmt.Errorf("statefulset %s is already duplicated", obj.Name)
	}

	// create a new StatefulSet and override the spec
	newName := fmt.Sprintf("%s-duplik8ted", statefulSet.Name)
	newStatefulSet := appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      newName,
			Namespace: statefulSet.Namespace,
			Labels: map[string]string{
				core.LABEL_DUPLICATED: "true",
			},
		},
		Spec: statefulSet.Spec,
	}

	// override the spec of the statefulset's pod
	configurator := clients.NewConfigurator(c.clientset, opts)
	err = configurator.OverrideSpec(c.ctx, obj.Namespace, &newStatefulSet.Spec.Template.Spec)
	if err != nil {
		return err
	}

	// create the new statefulset
	duplicatedStatefulSet, err := c.clientset.AppsV1().StatefulSets(obj.Namespace).Create(c.ctx, &newStatefulSet, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("statefulset %q duplicated in %q\n", obj.Name, newName)

	if opts.StartInteractiveShell {
		pod, err := GetOwnedPod(
			c.ctx,
			c.clientset,
			duplicatedStatefulSet.Namespace,
			duplicatedStatefulSet.Spec.Selector,
		)
		if err != nil {
			return err
		}
		return StartInteractiveShell(c.ctx, c.clientset, pod)
	}

	return nil
}
