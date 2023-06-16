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

package crud

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertHeaders(t *testing.T) {
	t.Run("convert from http.Headers to map", func(t *testing.T) {
		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("key1", "val1")

		clientOpts := ClientOptions{
			Headers: h,
		}

		res := clientOpts.convertHeaders()

		require.Equal(t, map[string]string{
			"Foo":  "bar",
			"Key1": "val1",
		}, res)
	})
}
