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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodClient struct {
	clientset *kubernetes.Clientset
	ctx       context.Context
}

func NewPodClient(opts utils.KubeOptions) (*PodClient, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &PodClient{
		clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

func (c *PodClient) Duplicate(obj core.DuplicableObject, opts core.DuplicateOpts) error {
	fmt.Printf("duplicating pod %s\n", obj.Name)

	// fetch the pod
	pod, err := c.clientset.CoreV1().Pods(obj.Namespace).Get(c.ctx, obj.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if pod.Labels[core.LABEL_DUPLICATED] == "true" {
		return fmt.Errorf("pod %s is already duplicated", obj.Name)
	}

	// create a new pod and override the spec
	newName := fmt.Sprintf("%s-duplik8ted", pod.Name)
	newPod := v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      newName,
			Namespace: pod.Namespace,
			Labels: map[string]string{
				core.LABEL_DUPLICATED: "true",
			},
		},
		Spec: pod.Spec,
	}

	// override the pod spec
	configurator := clients.NewConfigurator(c.clientset, opts)
	err = configurator.OverrideSpec(c.ctx, obj.Namespace, &newPod.Spec)
	if err != nil {
		return err
	}

	// create the new pod
	_, err = c.clientset.CoreV1().Pods(pod.Namespace).Create(c.ctx, &newPod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("pod %q duplicated in %q\n", obj.Name, newName)

	if opts.StartInteractiveShell {
		return StartInteractiveShell(c.ctx, c.clientset, newPod)
	}

	return nil
}
