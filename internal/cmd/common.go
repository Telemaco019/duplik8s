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
	"github.com/spf13/cobra"
	"github.com/telemaco019/duplik8s/internal/cmd/flags"
	"github.com/telemaco019/duplik8s/internal/utils"
)

func NewKubeOptions(cmd *cobra.Command, _ []string) (utils.KubeOptions, error) {
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
