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

package core

import "k8s.io/apimachinery/pkg/runtime/schema"

// ResourceSpec pairs a resource Kind with its GroupVersionResource.
// It is the single source of truth for the resources duplik8s supports.
type ResourceSpec struct {
	Kind string
	GVR  schema.GroupVersionResource
}

// SupportedResources is the list of resources duplik8s can duplicate.
// Commands, ListDuplicated and Delete all consume this list so that adding
// a new resource only requires extending it (plus a resource handler).
var SupportedResources = []ResourceSpec{
	{Kind: "Pod", GVR: schema.GroupVersionResource{Version: "v1", Resource: "pods"}},
	{Kind: "Deployment", GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}},
	{Kind: "StatefulSet", GVR: schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}},
}

// GVRForKind returns the GroupVersionResource for the given Kind, or false
// if the Kind is not supported.
func GVRForKind(kind string) (schema.GroupVersionResource, bool) {
	for _, r := range SupportedResources {
		if r.Kind == kind {
			return r.GVR, true
		}
	}
	return schema.GroupVersionResource{}, false
}
