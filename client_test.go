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
	"context"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("create new client", func(t *testing.T) {
		client, err := New[TestResource](ClientOptions{
			BaseURL: "http://crud-service/resource-path/",
		})
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("throws creating a new client", func(t *testing.T) {
		_, err := New[TestResource](ClientOptions{
			BaseURL: "http://crud-service/resource-path",
		})
		require.EqualError(t, err, fmt.Sprintf("%s: BaseURL must end with a trailing slash", ErrCreateClient))
	})
}

type TestResource struct {
	Field    string `json:"field"`
	ID       string `json:"_id"`
	IntField int    `json:"intField"`
}

func TestExport(t *testing.T) {
	ctx := context.Background()

	client, err := New[TestResource](ClientOptions{
		BaseURL: "http://crud-service/resource-path/",
	})
	require.NoError(t, err)

	t.Run("export data", func(t *testing.T) {
		responseBody := `
		{"field": "v-1","intField":1,"_id":"my-id-1"}
		{"field": "v-2","intField":2,"_id":"my-id-2"}
		{"field": "v-3","intField":3,"_id":"my-id-3"}
		`

		gock.New("http://crud-service/resource-path/").
			Get("export").
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, "", Filter{})
		require.NoError(t, err)
		require.Equal(t, []TestResource{
			{
				Field:    "v-1",
				IntField: 1,
				ID:       "my-id-1",
			},
			{
				Field:    "v-2",
				IntField: 2,
				ID:       "my-id-2",
			},
			{
				Field:    "v-3",
				IntField: 3,
				ID:       "my-id-3",
			},
		}, resources)
	})

	t.Run("export data with filter", func(t *testing.T) {
		responseBody := `
		{"field": "v-1","_id":"my-id-1"}
		{"field": "v-2","_id":"my-id-2"}
		`

		filter := Filter{
			Projection: []string{"field"},
			MongoQuery: map[string]any{
				"field": map[string]any{
					"$in": []string{"v-1", "v-2"},
				},
			},
			Limit: 5,
		}

		gock.New("http://crud-service/resource-path/").
			Get("export").
			AddMatcher(CrudQueryMatcher(t, filter)).
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, "", filter)
		require.NoError(t, err)
		require.Equal(t, []TestResource{
			{
				Field: "v-1",
				ID:    "my-id-1",
			},
			{
				Field: "v-2",
				ID:    "my-id-2",
			},
		}, resources)
	})

	t.Run("export data with only mongo query", func(t *testing.T) {
		responseBody := `
		{"field": "v-1","_id":"my-id-1"}
		{"field": "v-2","_id":"my-id-2"}
		`

		filter := Filter{
			MongoQuery: map[string]any{
				"field": map[string]any{"$in": []string{"v-1", "v-2"}},
			},
		}

		gock.New("http://crud-service/resource-path/").
			Get("export").
			AddMatcher(CrudQueryMatcher(t, filter)).
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, "", filter)
		require.NoError(t, err)
		require.Equal(t, []TestResource{
			{
				Field: "v-1",
				ID:    "my-id-1",
			},
			{
				Field: "v-2",
				ID:    "my-id-2",
			},
		}, resources)
	})

	t.Run("throws with errors", func(t *testing.T) {
		gock.New("http://crud-service/resource-path/").
			Get("export").
			Reply(500).
			BodyString(`{"message":"error message"}`)

		resources, err := client.Export(ctx, "", Filter{})
		require.EqualError(t, err, "GET http://crud-service/resource-path/export: 500 - {\"message\":\"error message\"}")
		require.Nil(t, resources)
	})
}
