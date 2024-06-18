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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodConfigurator struct {
	clientset *kubernetes.Clientset
	options   PodOverrideOptions
}

func NewPodConfigurator(clientset *kubernetes.Clientset, options PodOverrideOptions) PodConfigurator {
	return PodConfigurator{
		clientset: clientset,
		options:   options,
	}
}

func (c PodConfigurator) OverrideSpec(ctx context.Context, pod *v1.Pod) error {
	// Override command
	if c.options.Command != nil {
		for i := range pod.Spec.Containers {
			pod.Spec.Containers[i].Command = c.options.Command
			pod.Spec.Containers[i].Args = c.options.Args
			pod.Spec.Containers[i].ReadinessProbe = c.options.ReadinessProbe
			pod.Spec.Containers[i].LivenessProbe = c.options.LivenessProbe
			pod.Spec.Containers[i].ReadinessProbe = c.options.ReadinessProbe
			pod.Spec.Containers[i].StartupProbe = c.options.StartupProbe
		}
	}

	hasMountOncePvc, err := c.hasMountOncePvc(ctx, *pod)
	if err != nil {
		return err
	}

	// If the Pod does not have any PVC with mount once policy, then remove the node name
	// to allow the scheduler to schedule the pod on any node
	if !hasMountOncePvc {
		pod.Spec.NodeName = ""
	}

	// Override restart policy
	pod.Spec.RestartPolicy = v1.RestartPolicyNever

	return nil
}

func (c PodConfigurator) hasMountOncePvc(ctx context.Context, pod v1.Pod) (bool, error) {
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvc, err := c.clientset.
				CoreV1().
				PersistentVolumeClaims(pod.Namespace).
				Get(ctx, volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			return anyMountOnceAccessMode(pvc.Spec.AccessModes), nil
		}
	}
	return false, nil
}

func anyMountOnceAccessMode(modes []v1.PersistentVolumeAccessMode) bool {
	for _, mode := range modes {
		if mode == v1.ReadWriteOnce {
			return true
		}
		if mode == v1.ReadWriteOncePod {
			return true
		}
	}
	return false
}
