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
	"strings"

	"github.com/telemaco019/duplik8s/internal/clients"
	"github.com/telemaco019/duplik8s/internal/core"
	"github.com/telemaco019/duplik8s/internal/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

// resourceHandler captures the per-resource differences in the duplication
// flow. The Duplicator runs the shared flow and delegates the typed bits to
// these closures. Closures capture the clientset they need.
type resourceHandler struct {
	// kind is the resource Kind as used in TypeMeta (e.g. "Pod"). It is also
	// lowercased for user-facing messages.
	kind       string
	apiVersion string
	fetch      func(ctx context.Context, namespace, name string) (runtime.Object, error)
	podSpec    func(obj runtime.Object) *v1.PodSpec
	buildNew   func(original runtime.Object) runtime.Object
	create     func(ctx context.Context, obj runtime.Object) (runtime.Object, error)
	deleteObj  func(ctx context.Context, obj runtime.Object) error
	// ownedPod resolves the Pod to exec into after a controller is duplicated.
	// When nil, the created object is itself a Pod and is used directly.
	ownedPod func(ctx context.Context, clientset *kubernetes.Clientset, created runtime.Object) (v1.Pod, error)
}

// Duplicator is the unified implementation of core.Duplicator. It holds the
// clientset and a resourceHandler describing how to duplicate a given kind.
type Duplicator struct {
	clientset *kubernetes.Clientset
	handler   resourceHandler
}

func newDuplicator(opts utils.KubeOptions, handlerFactory func(*kubernetes.Clientset) resourceHandler) (*Duplicator, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &Duplicator{
		clientset: clientset,
		handler:   handlerFactory(clientset),
	}, nil
}

func (d *Duplicator) Duplicate(ctx context.Context, obj core.DuplicableObject, opts core.DuplicateOpts) error {
	kind := strings.ToLower(d.handler.kind)
	fmt.Printf("duplicating %s %s\n", kind, obj.Name)

	// fetch the resource
	original, err := d.handler.fetch(ctx, obj.Namespace, obj.Name)
	if err != nil {
		return err
	}
	accessor, err := meta.Accessor(original)
	if err != nil {
		return fmt.Errorf("cannot read labels from %s %s: %w", kind, obj.Name, err)
	}
	if accessor.GetLabels()[core.LABEL_DUPLICATED] == "true" {
		return fmt.Errorf("%s %s is already duplicated", kind, obj.Name)
	}

	// build the new object with the duplicated name/label and a copied Spec
	newObj := d.handler.buildNew(original)

	// set the TypeMeta centrally from the handler's kind/apiVersion
	gv, _ := schema.ParseGroupVersion(d.handler.apiVersion)
	newObj.GetObjectKind().SetGroupVersionKind(gv.WithKind(d.handler.kind))

	// override the pod spec
	configurator := clients.NewConfigurator(d.clientset, opts)
	if err = configurator.OverrideSpec(ctx, obj.Namespace, d.handler.podSpec(newObj)); err != nil {
		return err
	}

	// create the new resource
	created, err := d.handler.create(ctx, newObj)
	if err != nil {
		return err
	}
	newName := fmt.Sprintf("%s-duplik8ted", obj.Name)
	fmt.Printf("%s %q duplicated in %q\n", kind, obj.Name, newName)

	if !opts.StartInteractiveShell {
		return nil
	}

	// resolve the pod to exec into
	var pod v1.Pod
	if d.handler.ownedPod != nil {
		pod, err = d.handler.ownedPod(ctx, d.clientset, created)
		if err != nil {
			return err
		}
	} else {
		// the created object is itself the pod
		pod = *created.(*v1.Pod)
	}
	return StartInteractiveShell(ctx, d.clientset, pod, created, d.handler.deleteObj)
}
