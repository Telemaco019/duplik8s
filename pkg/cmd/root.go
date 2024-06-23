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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/pkg/core"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func NewRootCmd(
	podClient core.Duplik8sClient,
	deployClient core.Duplik8sClient,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "kubectl-duplicate",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl duplicate",
		},
		Short: "duplik8s is a kubectl plugin for duplicating Kubernetes resources.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Setup kubeconfig flags
	defaultNamespace := "default"
	defaultKubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfig = filepath.Join(home, ".kube", "config")
	}
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.KubeConfig = &defaultKubeconfig
	configFlags.Namespace = &defaultNamespace
	configFlags.AddFlags(rootCmd.PersistentFlags())

	// add subcommands
	rootCmd.AddCommand(NewPodCmd(podClient))
	rootCmd.AddCommand(NewDeployCmd(deployClient))

	return rootCmd
}
