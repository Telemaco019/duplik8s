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
	"github.com/telemaco019/duplik8s/pkg/utils"
)

type duplicatorFactory func(opts utils.KubeOptions) (core.Duplicator, error)

func newDuplicateCmd(factory duplicatorFactory, client core.Client, selectMessage string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts, err := NewKubeOptions(cmd, args)
		if err != nil {
			return err
		}
		duplicator, err := factory(opts)
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
			obj, err = utils.SelectItem(client, opts.Namespace, selectMessage)
			if err != nil {
				return err
			}
		} else {
			obj = core.DuplicableObject{
				Name:      args[0],
				Namespace: opts.Namespace,
			}
		}

		return duplicator.Duplicate(obj, options)
	}
}

func addOverrideFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice(
		flags.COMMAND_OVERRIDE,
		[]string{"/bin/sh"},
		"Override the command of each container in the Pod.",
	)
	cmd.Flags().StringSlice(
		flags.ARGS_OVERRIDE,
		[]string{"-c", "trap 'exit 0' INT TERM KILL; while true; do sleep 1; done"},
		"Override the command of each container in the Pod.",
	)
}
