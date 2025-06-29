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

package utils

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

// NewClientset creates a new kubernetes clientset
func NewClientset(kubeconfig, context string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			ClusterInfo:    clientcmdapi.Cluster{Server: ""},
			CurrentContext: context,
		},
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func getKubeClientConfig(kubeconfig, context string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			ClusterInfo:    clientcmdapi.Cluster{Server: ""},
			CurrentContext: context,
		},
	).ClientConfig()
}

func NewDynamicClient(kubeconfig, context string) (*dynamic.DynamicClient, error) {
	config, err := getKubeClientConfig(kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientSet, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func GetKubeconfigPath() string {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath != "" {
		return kubeconfigPath
	}
	home := homedir.HomeDir()
	return filepath.Join(home, ".kube", "config")
}

func NewDiscoveryClient(kubeconfig, context string) (*discovery.DiscoveryClient, error) {
	config, err := getKubeClientConfig(kubeconfig, context)
	if err != nil {
		return nil, err
	}

	clientSet, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func IsPodReady(pod *v1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func WaitUntilPodReady(ctx context.Context, client kubernetes.Interface, pod v1.Pod, timeout time.Duration) error {
	watchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	watch, err := client.CoreV1().Pods(pod.Namespace).Watch(watchCtx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", pod.Name),
	})
	if err != nil {
		return err
	}
	defer watch.Stop()

	for event := range watch.ResultChan() {
		if pod, ok := event.Object.(*v1.Pod); ok && IsPodReady(pod) {
			return nil
		}
	}
	return fmt.Errorf("pod %s not ready within timeout", pod.Name)
}

type KubeOptions struct {
	Kubeconfig  string
	Kubecontext string
	Namespace   string
}
