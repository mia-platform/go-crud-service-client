// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testhelper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func ParseResponseToNdjson[TResource any](t *testing.T, response []TResource) string {
	t.Helper()

	var responseBytes = []byte("")
	separator := []byte("\n")

	for _, elem := range response {
		marshElem, err := json.Marshal(elem)
		require.NoError(t, err, "Unexpected error to pars elem in NdJSON")

		responseBytes = append(responseBytes, separator...)
		responseBytes = append(responseBytes, marshElem...)
	}

	return string(responseBytes)
}
