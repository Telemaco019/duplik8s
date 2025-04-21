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
	"fmt"
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
	cmd := NewRootCmd(podClient, podClient)
	output, err := test.ExecuteCommand(cmd, "pod")
	assert.NotEmpty(t, output)
	assert.Error(t, err)
}

func Test_Success(t *testing.T) {
	podClient := mocks.NewPodClient(
		mocks.ListPodsResult{},
		nil,
	)
	cmd := NewRootCmd(podClient, podClient)
	_, err := test.ExecuteCommand(cmd, "pod", "pod-1")
	assert.NoError(t, err)
}

func Test_DuplicateError(t *testing.T) {
	podClient := mocks.NewPodClient(
		mocks.ListPodsResult{},
		fmt.Errorf("error"),
	)
	cmd := NewRootCmd(podClient, podClient)
	_, err := test.ExecuteCommand(cmd, "pod", "pod-1")
	assert.EqualError(t, err, "error")
}
