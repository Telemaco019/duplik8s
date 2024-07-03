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

package utils

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/telemaco019/duplik8s/pkg/core"
)

func SelectItem(client core.Client, namespace, selectMessage string) (core.DuplicableObject, error) {
	var selected = core.DuplicableObject{}
	objs, err := client.ListDuplicable(namespace)
	if err != nil {
		return selected, err
	}
	if len(objs) == 0 {
		return selected, fmt.Errorf("no Pods found in namespace %q", namespace)
	}
	options := make([]huh.Option[core.DuplicableObject], len(objs))
	for i, o := range objs {
		options[i] = huh.NewOption(o.Name, o)
	}
	err = huh.NewSelect[core.DuplicableObject]().
		Title(fmt.Sprintf("%s [%s]", selectMessage, namespace)).
		Options(options...).
		Value(&selected).
		Run()
	return selected, err
}
