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
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/pkg/clients"
	"github.com/telemaco019/duplik8s/pkg/core"
)

func cleanup(client core.Client, namespace string) error {
	duplicated, err := client.ListDuplicated(context.Background(), namespace)
	if err != nil {
		return err
	}
	if len(duplicated) == 0 {
		fmt.Printf("No duplicated resources found in namespace %q\n", namespace)
		return nil
	}

	renderDuplicatedObjects(duplicated)

	var shouldDelete bool
	err = huh.NewConfirm().Title("Do you want to delete the following resources?").Value(&shouldDelete).Run()
	if err != nil {
		return err
	}

	if shouldDelete {
		for _, obj := range duplicated {
			err = client.Delete(context.Background(), obj)
			if err != nil {
				return err
			}
			fmt.Printf("deleted %s %s/%s\n", obj.ObjectKind.GroupVersionKind().Kind, obj.Namespace, obj.Name)
		}
	}

	return nil
}

func NewCleanupCmd(client core.Client) *cobra.Command {
	podCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup duplicated resources.",
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
			return cleanup(client, opts.Namespace)
		},
	}
	addOverrideFlags(podCmd)
	return podCmd
}
