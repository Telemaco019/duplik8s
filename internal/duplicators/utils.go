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
	"github.com/charmbracelet/huh"
	"github.com/telemaco019/duplik8s/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"time"
)

func StartInteractiveShell(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	pod corev1.Pod,
	duplicatedObject runtime.Object,
) error {
	// wait for the pod to be ready
	fmt.Printf("waiting for the duplicated pod %q to be ready...\n", pod.Name)
	err := utils.WaitUntilPodReady(ctx, clientset, pod, 60*time.Second)
	if err != nil {
		return err
	}
	fmt.Println("Pod is ready, launching shell...")
	execCmd := []string{
		"kubectl", "exec", "-it", pod.Name, "-n", pod.Namespace, "--", "/bin/sh",
	}
	if err = utils.RunInteractive(execCmd); err != nil {
		return fmt.Errorf("error during shell session: %w", err)
	}

	// prompt for deletion
	var confirmDelete = true
	err = huh.NewConfirm().Title(
		"Do you want to delete the duplicated resource?",
	).Value(&confirmDelete).Run()
	if err != nil {
		return err
	}
	if confirmDelete {
		// Delete the duplicated object based on its type
		err = deleteResource(ctx, clientset, duplicatedObject)
		if err != nil {
			return err
		}
		fmt.Println("duplicated resource deleted.")
	} else {
		fmt.Println("duplicated resource retained.")
	}

	return nil
}

func deleteResource(ctx context.Context, clientset *kubernetes.Clientset, duplicatedObject runtime.Object) error {
	switch obj := duplicatedObject.(type) {
	case *corev1.Pod:
		return clientset.CoreV1().Pods(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
	case *appsv1.Deployment:
		return clientset.AppsV1().Deployments(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
	case *appsv1.StatefulSet:
		return clientset.AppsV1().StatefulSets(obj.Namespace).Delete(ctx, obj.Name, metav1.DeleteOptions{})
	default:
		return fmt.Errorf("unsupported duplicated object type: %T", obj)
	}
}

func GetOwnedPod(
	ctx context.Context,
	client *kubernetes.Clientset,
	namespace string,
	selector *metav1.LabelSelector,
) (corev1.Pod, error) {
	labelSelector := metav1.FormatLabelSelector(selector)

	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return corev1.Pod{}, fmt.Errorf("failed to list pods with selector %q: %w", labelSelector, err)
	}

	if len(podList.Items) == 0 {
		return corev1.Pod{}, fmt.Errorf("no pods found with selector %q", labelSelector)
	}

	return podList.Items[0], nil
}
