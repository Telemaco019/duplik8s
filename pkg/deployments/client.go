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

package deployments

import (
	"context"
	"fmt"
	"github.com/telemaco019/duplik8s/pkg/core"
	"github.com/telemaco019/duplik8s/pkg/pods"
	"github.com/telemaco019/duplik8s/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentClient struct {
	clientset *kubernetes.Clientset
	ctx       context.Context
}

func NewClient(opts utils.KubeOptions) (*DeploymentClient, error) {
	clientset, err := utils.NewClientset(opts.Kubeconfig, opts.Kubecontext)
	if err != nil {
		return nil, err
	}
	return &DeploymentClient{
		clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

func (c *DeploymentClient) List(namespace string) ([]core.DuplicableObject, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(c.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var objs []core.DuplicableObject
	for _, d := range deployments.Items {
		objs = append(objs, core.NewDeployment(d.Name, d.Namespace))
	}
	return objs, nil
}

func (c *DeploymentClient) Duplicate(obj core.DuplicableObject, opts core.PodOverrideOptions) error {
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
	configurator := pods.NewConfigurator(c.clientset, opts)
	err = configurator.OverrideSpec(c.ctx, obj.Namespace, &newDeploy.Spec.Template.Spec)
	if err != nil {
		return err
	}

	// create the new deployment
	_, err = c.clientset.AppsV1().Deployments(obj.Namespace).Create(c.ctx, &newDeploy, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("deployment %q duplicated in %q\n", obj.Name, newName)
	return nil
}
