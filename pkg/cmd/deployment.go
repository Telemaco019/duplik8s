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
	"github.com/telemaco019/duplik8s/pkg/cmd/flags"
	"github.com/telemaco019/duplik8s/pkg/core"
	"github.com/telemaco019/duplik8s/pkg/deployments"
	"github.com/telemaco019/duplik8s/pkg/utils"
)

func NewDeployCmd(client *deployments.DeploymentClient) *cobra.Command {
	podCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Duplicate a Deployment.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := NewKubeOptions(cmd, args)
			if err != nil {
				return err
			}
			if client == nil {
				client, err = deployments.NewClient(opts)
				if err != nil {
					return err
				}
			}
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
			options := core.PodOverrideOptions{
				Command: cmdOverride,
				Args:    argsOverride,
			}

			var obj core.DuplicableObject
			if len(args) == 0 {
				obj, err = utils.SelectItem(client, opts.Namespace, "Select a Deployment")
				if err != nil {
					return err
				}
			} else {
				obj = core.NewPod(args[0], opts.Namespace)
			}

			return client.Duplicate(obj, options)
		},
	}
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

	return podCmd
}
