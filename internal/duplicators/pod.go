/*
 * Copyright 2026 Michele Zanotti <m.zanotti019@gmail.com>
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

	"github.com/telemaco019/duplik8s/internal/core"
	"github.com/telemaco019/duplik8s/internal/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// NewPodClient builds a Duplicator for Pods.
func NewPodClient(opts utils.KubeOptions) (*Duplicator, error) {
	return newDuplicator(opts, newPodHandler)
}

func newPodHandler(clientset *kubernetes.Clientset) resourceHandler {
	return resourceHandler{
		kind:       "Pod",
		apiVersion: "v1",
		fetch: func(ctx context.Context, namespace, name string) (runtime.Object, error) {
			return clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		},
		podSpec: func(obj runtime.Object) *v1.PodSpec {
			return &obj.(*v1.Pod).Spec
		},
		buildNew: func(original runtime.Object) runtime.Object {
			pod := original.(*v1.Pod)
			return &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-duplik8ted", pod.Name),
					Namespace: pod.Namespace,
					Labels: map[string]string{
						core.LABEL_DUPLICATED: "true",
					},
				},
				Spec: pod.Spec,
			}
		},
		create: func(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
			newPod := obj.(*v1.Pod)
			return clientset.CoreV1().Pods(newPod.Namespace).Create(ctx, newPod, metav1.CreateOptions{})
		},
		deleteObj: func(ctx context.Context, obj runtime.Object) error {
			pod := obj.(*v1.Pod)
			return clientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
		},
		// ownedPod is nil: the created object is itself the Pod.
		ownedPod: nil,
	}
}
