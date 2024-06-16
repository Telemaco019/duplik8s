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
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/pkg/cmd/flags"
	"github.com/telemaco019/duplik8s/pkg/pods"
	"github.com/telemaco019/duplik8s/pkg/utils"
)

func NewKubeOptions(cmd *cobra.Command, args []string) (utils.KubeOptions, error) {
	var err error

	o := utils.KubeOptions{}
	o.Kubeconfig, err = cmd.Flags().GetString(flags.KUBECONFIG)
	if err != nil {
		return o, err
	}
	o.Kubecontext, err = cmd.Flags().GetString(flags.KUBECONTEXT)
	if err != nil {
		return o, err
	}
	o.Namespace, err = cmd.Flags().GetString(flags.NAMESPACE)
	if err != nil {
		return o, err
	}

	return o, nil
}

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Duplicate a Pod.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opts, err := NewKubeOptions(cmd, args)
		if err != nil {
			return err
		}
		client, err := pods.NewClient(opts)
		if err != nil {
			return err
		}
		cmdOverride, err := cmd.Flags().GetStringSlice(flags.COMMAND_OVERRIDE)
		if err != nil {
			return err
		}
		argsOverride, err := cmd.Flags().GetStringSlice(flags.ARGS_OVERRIDE)
		if err != nil {
			return err
		}

		// Avoid printing usage information on errors
		cmd.SilenceUsage = true
		options := pods.PodOverrideOptions{
			Command: cmdOverride,
			Args:    argsOverride,
		}

		var podName string
		if len(args) == 0 {
			podName, err = selectPod(client, opts.Namespace)
			if err != nil {
				return err
			}
		} else {
			podName = args[0]
		}

		return client.DuplicatePod(podName, opts.Namespace, options)
	},
}

func selectPod(client *pods.PodClient, namespace string) (string, error) {
	availablePods, err := client.ListPods(namespace)
	if err != nil {
		return "", err
	}
	if len(availablePods) == 0 {
		return "", fmt.Errorf("no Pods found in namespace %q", namespace)
	}
	options := make([]huh.Option[string], len(availablePods))
	for i, p := range availablePods {
		options[i] = huh.NewOption(p, p)
	}
	var selectedPod string
	err = huh.NewSelect[string]().
		Title(fmt.Sprintf("Select a Pod [%s]", namespace)).
		Options(options...).
		Value(&selectedPod).
		Run()
	return selectedPod, err
}

func init() {
	rootCmd.AddCommand(podCmd)

	podCmd.Flags().StringSlice(
		flags.COMMAND_OVERRIDE,
		[]string{"/bin/sh"},
		"Override the command of each container in the Pod.",
	)
	podCmd.Flags().StringSlice(
		flags.ARGS_OVERRIDE,
		[]string{"-c", "trap 'exit 0' INT TERM KILL; while true; do sleep 1; done"},
		"Override the command of each container in the Pod.",
	)
}
