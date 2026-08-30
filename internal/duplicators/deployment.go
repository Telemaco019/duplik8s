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

	"github.com/telemaco019/duplik8s/internal/core"
	"github.com/telemaco019/duplik8s/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// NewDeploymentClient builds a Duplicator for Deployments.
func NewDeploymentClient(opts utils.KubeOptions) (*Duplicator, error) {
	return newDuplicator(opts, newDeploymentHandler)
}

func newDeploymentHandler(clientset *kubernetes.Clientset) resourceHandler {
	return resourceHandler{
		kind:       "Deployment",
		apiVersion: "v1",
		fetch: func(ctx context.Context, namespace, name string) (runtime.Object, error) {
			return clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		},
		podSpec: func(obj runtime.Object) *v1.PodSpec {
			return &obj.(*appsv1.Deployment).Spec.Template.Spec
		},
		buildNew: func(original runtime.Object) runtime.Object {
			deploy := original.(*appsv1.Deployment)
			return &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-duplik8ted", deploy.Name),
					Namespace: deploy.Namespace,
					Labels: map[string]string{
						core.LABEL_DUPLICATED: "true",
					},
				},
				Spec: deploy.Spec,
			}
		},
		create: func(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
			newDeploy := obj.(*appsv1.Deployment)
			return clientset.AppsV1().Deployments(newDeploy.Namespace).Create(ctx, newDeploy, metav1.CreateOptions{})
		},
		deleteObj: func(ctx context.Context, obj runtime.Object) error {
			deploy := obj.(*appsv1.Deployment)
			return clientset.AppsV1().Deployments(deploy.Namespace).Delete(ctx, deploy.Name, metav1.DeleteOptions{})
		},
		ownedPod: func(ctx context.Context, cs *kubernetes.Clientset, created runtime.Object) (v1.Pod, error) {
			deploy := created.(*appsv1.Deployment)
			return GetOwnedPod(ctx, cs, deploy.Namespace, deploy.Spec.Selector)
		},
	}
}
