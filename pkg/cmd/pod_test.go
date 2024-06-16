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
	"github.com/stretchr/testify/assert"
	"github.com/telemaco019/duplik8s/pkg/test"
	"github.com/telemaco019/duplik8s/pkg/test/mocks"
	"testing"
)

func Test_NoPodsAvailable(t *testing.T) {
	podClient := mocks.NewPodClient(
		mocks.ListPodsResult{},
		nil,
	)
	cmd := NewRootCmd(podClient)
	output, err := test.ExecuteCommand(cmd, "pod")
	assert.Equal(t, "", output)
	assert.Error(t, err)
}