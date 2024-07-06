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
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/pkg/clients"
	"github.com/telemaco019/duplik8s/pkg/cmd/flags"
	"github.com/telemaco019/duplik8s/pkg/core"
	"github.com/telemaco019/duplik8s/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type duplicatorFactory func(opts utils.KubeOptions) (core.Duplicator, error)

func newDuplicateCmd(newDuplicator duplicatorFactory, client core.Client, gvr schema.GroupVersionResource) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts, err := NewKubeOptions(cmd, args)
		if err != nil {
			return err
		}
		duplicator, err := newDuplicator(opts)
		if err != nil {
			return err
		}
		if client == nil {
			client, err = clients.NewDuplik8sClient(opts)
			if err != nil {
				return err
			}
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

		// If available, duplicate the resource provided as argument
		var obj core.DuplicableObject
		if len(args) > 0 {
			obj = core.DuplicableObject{
				Name:      args[0],
				Namespace: opts.Namespace,
			}
			return duplicator.Duplicate(obj, options)
		}

		// Otherwise, list available resources
		objs, err := client.ListDuplicable(
			context.Background(),
			gvr,
			opts.Namespace,
		)
		if err != nil {
			return err
		}
		if len(objs) == 0 {
			return fmt.Errorf("no %s available in namespace %q", gvr.Resource, opts.Namespace)
		}
		caser := cases.Title(language.English)
		obj, err = utils.SelectItem(
			objs,
			fmt.Sprintf("%s [%s]", caser.String(gvr.Resource), opts.Namespace),
		)
		if err != nil {
			return err
		}
		obj = core.DuplicableObject{
			Name:      args[0],
			Namespace: opts.Namespace,
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
