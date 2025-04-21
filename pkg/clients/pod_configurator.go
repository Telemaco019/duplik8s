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
	"github.com/telemaco019/duplik8s/pkg/core"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodConfigurator struct {
	clientset *kubernetes.Clientset
	options   core.PodOverrideOptions
}

func NewConfigurator(clientset *kubernetes.Clientset, options core.PodOverrideOptions) PodConfigurator {
	return PodConfigurator{
		clientset: clientset,
		options:   options,
	}
}

func (c PodConfigurator) OverrideSpec(ctx context.Context, namespace string, podSpec *v1.PodSpec) error {
	// Override command
	if c.options.Command != nil {
		for i := range podSpec.Containers {
			podSpec.Containers[i].Command = c.options.Command
			podSpec.Containers[i].Args = c.options.Args
			podSpec.Containers[i].ReadinessProbe = c.options.ReadinessProbe
			podSpec.Containers[i].LivenessProbe = c.options.LivenessProbe
			podSpec.Containers[i].ReadinessProbe = c.options.ReadinessProbe
			podSpec.Containers[i].StartupProbe = c.options.StartupProbe
		}
	}

	hasMountOncePvc, err := c.hasMountOncePvc(ctx, namespace, *podSpec)
	if err != nil {
		return err
	}

	// If the Pod does not have any PVC with mount once policy, then remove the node name
	// to allow the scheduler to schedule the pod on any node
	if !hasMountOncePvc {
		podSpec.NodeName = ""
	}

	return nil
}

func (c PodConfigurator) hasMountOncePvc(ctx context.Context, namespace string, podSpec v1.PodSpec) (bool, error) {
	for _, volume := range podSpec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvc, err := c.clientset.
				CoreV1().
				PersistentVolumeClaims(namespace).
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
