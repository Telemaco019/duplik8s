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

package core

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Duplicator interface {
	Duplicate(obj DuplicableObject, opts PodOverrideOptions) error
}

type Client interface {
	ListDuplicable(
		ctx context.Context,
		resource schema.GroupVersionResource,
		namespace string,
	) ([]DuplicableObject, error)
	ListDuplicated(ctx context.Context, namespace string) ([]DuplicatedObject, error)
	Delete(ctx context.Context, obj DuplicatedObject) error
}

type PodOverrideOptions struct {
	// Command overrides the default command of each container.
	Command []string
	// Args overrides the default args of each container.
	Args []string
	// ReadinessProbe overrides the readiness probe of each container.
	ReadinessProbe *v1.Probe
	// LivenessProbe overrides the liveness probe of each container.
	LivenessProbe *v1.Probe
	// StartupProbe overrides the startup probe of each container.
	StartupProbe *v1.Probe
}

type DuplicatedObject struct {
	Name              string
	Namespace         string
	ObjectKind        schema.ObjectKind
	CreationTimestamp metav1.Time
}

type DuplicableObject struct {
	Name      string
	Namespace string
}

func NewDuplicable(u unstructured.Unstructured) DuplicableObject {
	return DuplicableObject{
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
	}
}
