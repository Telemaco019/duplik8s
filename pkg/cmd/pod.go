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
	"github.com/telemaco019/duplik8s/pkg/duplicators"
	"github.com/telemaco019/duplik8s/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewPodCmd(duplicator core.Duplicator, client core.Client) *cobra.Command {
	factory := func(opts utils.KubeOptions) (core.Duplicator, error) {
		if duplicator == nil {
			return duplicators.NewPodClient(opts)
		}
		return duplicator, nil
	}
	podCmd := &cobra.Command{
		Use:   "pod",
		Short: "Duplicate a Pod.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			run := newDuplicateCmd(factory, client, schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			})
			return run(cmd, args)
		},
	}
	addOverrideFlags(podCmd)
	return podCmd
}
