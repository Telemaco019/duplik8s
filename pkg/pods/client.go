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

package pods

import (
	"context"
	"fmt"
	"github.com/telemaco019/duplik8s/pkg/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodClient struct {
	clientset *kubernetes.Clientset
	ctx       context.Context
}

func NewClient(opts utils.KubeOptions) (*PodClient, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &PodClient{
		clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

type PodOverrideOptions struct {
	// Command overrides the default command of each container.
	Command []string
}

func (o PodOverrideOptions) Apply(pod *v1.Pod) {
	// Override command
	if o.Command != nil {
		for i := range pod.Spec.Containers {
			pod.Spec.Containers[i].Command = o.Command
		}
	}
}

func (c PodClient) DuplicatePod(podName string, namespace string, opts PodOverrideOptions) error {
	fmt.Printf("duplicating Pod %s\n", podName)

	// fetch the pod
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(c.ctx, podName, metav1.GetOptions{})
	if err != nil {
		return err
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
		},
		Spec: pod.Spec,
	}
	opts.Apply(&newPod)

	// create the new pod
	_, err = c.clientset.CoreV1().Pods(pod.Namespace).Create(c.ctx, &newPod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Pod %s duplicated in %s\n", podName, newName)
	return nil
}
