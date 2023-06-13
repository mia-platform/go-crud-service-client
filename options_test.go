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
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mia-platform/go-crud-service-client/internal/types"

	"github.com/stretchr/testify/require"
)

func TestSetOptionsInRequest(t *testing.T) {
	t.Run("set both query and headers", func(t *testing.T) {
		h := http.Header{}
		h.Add("foo", "bar")
		h.Add("Authorization", "Bearer taz")

		options := Options{
			Filter: Filter{
				Limit: 5,
			},
			Headers: h,
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		err := options.setOptionsInRequest(req)

		require.NoError(t, err)
		require.Equal(t, h, req.Header)
		require.Equal(t, "5", req.URL.Query().Get("_l"))
	})
}

func TestAddCrudQueryToRequest(t *testing.T) {
	tests := []struct {
		name                   string
		filter                 types.Filter
		expectedUnencodedQuery string
	}{
		{
			name:                   "with empty filter returns empty query",
			filter:                 types.Filter{},
			expectedUnencodedQuery: "",
		},
		{
			name: "with all filters",
			filter: types.Filter{
				Limit:      5,
				Projection: []string{"a", "b"},
				FiledsQuery: map[string]string{
					"customId": "abcde",
				},
				MongoQuery: map[string]any{
					"field": map[string]any{
						"$in": []string{"v-1", "v-2"},
					},
				},
				Skip: 2,
			},
			expectedUnencodedQuery: `_l=5&_p=a,b&customId=abcde&_q={"field":{"$in":["v-1","v-2"]}}&_sk=2`,
		},
		{
			name: "with only FiledsQuery",
			filter: types.Filter{
				FiledsQuery: map[string]string{
					"customId": "abcde",
					"name":     "Alice",
				},
			},
			expectedUnencodedQuery: `customId=abcde&name=Alice`,
		},
		{
			name: "with only MongoQuery",
			filter: types.Filter{
				MongoQuery: map[string]any{
					"field": map[string]any{
						"$in": []string{"v-1", "v-2"},
					},
				},
			},
			expectedUnencodedQuery: `_q={"field":{"$in":["v-1","v-2"]}}`,
		},
		{
			name: "with only limit",
			filter: types.Filter{
				Limit: 5,
			},
			expectedUnencodedQuery: `_l=5`,
		},
		{
			name: "with only projection",
			filter: types.Filter{
				Projection: []string{"a", "b"},
			},
			expectedUnencodedQuery: `_p=a,b`,
		},
		{
			name: "with only skip",
			filter: types.Filter{
				Skip: 4,
			},
			expectedUnencodedQuery: `_sk=4`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			err := addCrudQueryToRequest(req, test.filter)

			require.NoError(t, err)
			q, err := url.ParseQuery(test.expectedUnencodedQuery)
			require.NoError(t, err)
			require.Equal(t, q.Encode(), req.URL.RawQuery)
		})
	}
}

func TestAddHeaderToRequest(t *testing.T) {
	t.Run("with nil headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		addHeaderToRequest(req, nil)

		require.Empty(t, req.Header)
	})

	t.Run("with empty headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		addHeaderToRequest(req, http.Header{})

		require.Empty(t, req.Header)
	})

	t.Run("with headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		h := http.Header{}
		h.Add("foo", "bar")
		h.Add("Authorization", "Bearer taz")

		addHeaderToRequest(req, h)

		require.Equal(t, h, req.Header)
	})
}
