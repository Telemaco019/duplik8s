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

type DeploymentClient struct {
	clientset *kubernetes.Clientset
	ctx       context.Context
}

func NewDeploymentClient(opts utils.KubeOptions) (*DeploymentClient, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &DeploymentClient{
		clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

func (c *DeploymentClient) Duplicate(obj core.DuplicableObject, opts core.DuplicateOpts) error {
	fmt.Printf("duplicating deployment %s\n", obj.Name)

	// fetch the Deployment
	deploy, err := c.clientset.AppsV1().Deployments(obj.Namespace).Get(c.ctx, obj.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if deploy.Labels[core.LABEL_DUPLICATED] == "true" {
		return fmt.Errorf("deployment %s is already duplicated", obj.Name)
	}

	// create a new Deployment and override the spec
	newName := fmt.Sprintf("%s-duplik8ted", deploy.Name)
	newDeploy := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      newName,
			Namespace: deploy.Namespace,
			Labels: map[string]string{
				core.LABEL_DUPLICATED: "true",
			},
		},
		Spec: deploy.Spec,
	}

	// override the spec of the deployment's pod
	configurator := clients.NewConfigurator(c.clientset, opts)
	err = configurator.OverrideSpec(c.ctx, obj.Namespace, &newDeploy.Spec.Template.Spec)
	if err != nil {
		return err
	}

	// create the new deployment
	duplicatedDeploy, err := c.clientset.AppsV1().Deployments(obj.Namespace).Create(c.ctx, &newDeploy, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("deployment %q duplicated in %q\n", obj.Name, newName)

	if opts.StartInteractiveShell {
		pod, err := GetOwnedPod(
			c.ctx,
			c.clientset,
			newDeploy.Namespace,
			duplicatedDeploy.Spec.Selector,
		)
		if err != nil {
			return err
		}
		return StartInteractiveShell(c.ctx, c.clientset, pod)
	}

	return nil
}
