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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func crudQueryMatcher(t *testing.T, expectedFilter Filter) gock.MatchFunc {
	return func(r1 *http.Request, r2 *gock.Request) (bool, error) {
		t.Helper()

		actualQuery := r1.URL.Query()

		if expectedFilter.MongoQuery != nil {
			actualMongoQuery := actualQuery.Get("_q")

			expectedQueryBytes, err := json.Marshal(expectedFilter.MongoQuery)
			require.NoError(t, err)

			if !assert.JSONEq(t, string(expectedQueryBytes), actualMongoQuery) {
				return false, fmt.Errorf("mongo query check fails. Actual: %s, required: %+v", actualMongoQuery, expectedFilter.MongoQuery)
			}
		}

		if expectedFilter.Projection != nil {
			actualProjection := actualQuery.Get("_p")

			if !assert.Equal(t, strings.Join(expectedFilter.Projection, ","), actualProjection) {
				return false, fmt.Errorf("projection query check fails. Actual: %s, required: %s", actualProjection, expectedFilter.Projection)
			}
		}

		if expectedFilter.Limit != 0 {
			actualLimit := actualQuery.Get("_l")

			if !assert.Equal(t, strconv.Itoa(expectedFilter.Limit), actualLimit) {
				return false, fmt.Errorf("limit query check fails. Actual: %s, required: %d", actualLimit, expectedFilter.Limit)
			}
		}

		return true, nil
	}
}
