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

package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/internal/clients"
	"github.com/telemaco019/duplik8s/internal/core"
)

func listDuplicatedResources(client core.Client, namespace string) error {
	duplicatedObjs, err := client.ListDuplicated(context.Background(), namespace)
	if err != nil {
		return err
	}

	if len(duplicatedObjs) == 0 {
		fmt.Printf("No duplicated resources found in namespace %q\n", namespace)
		return nil
	}
	renderDuplicatedObjects(duplicatedObjs)
	return nil
}

func NewListDuplicatedCmd(client core.Client) *cobra.Command {
	podCmd := &cobra.Command{
		Use:   "list",
		Short: "Show duplicated resources.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			opts, err := NewKubeOptions(cmd, args)
			if err != nil {
				return err
			}
			if client == nil {
				client, err = clients.NewDuplik8sClient(opts)
				if err != nil {
					return err
				}
			}
			return listDuplicatedResources(client, opts.Namespace)
		},
	}
	addOverrideFlags(podCmd)
	return podCmd
}
