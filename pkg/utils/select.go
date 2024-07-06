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
	"github.com/charmbracelet/huh"
	"github.com/telemaco019/duplik8s/pkg/core"
)

func SelectItem(items []core.DuplicableObject, selectMessage string) (core.DuplicableObject, error) {
	var selected core.DuplicableObject
	options := make([]huh.Option[core.DuplicableObject], len(items))
	for i, o := range items {
		options[i] = huh.NewOption(o.Name, o)
	}
	err := huh.NewSelect[core.DuplicableObject]().
		Title(selectMessage).
		Options(options...).
		Value(&selected).
		Run()
	return selected, err
}
