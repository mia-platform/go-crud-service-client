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
	"net/http"
	"testing"

	"github.com/mia-platform/go-crud-service-client/testhelper"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

var baseURL = "http://crud-service/resource-path/"

func TestNewClient(t *testing.T) {
	t.Run("create new client", func(t *testing.T) {
		client, err := NewClient[TestResource](ClientOptions{
			BaseURL: baseURL,
		})
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("create a new client correctly without trailing slash", func(t *testing.T) {
		c, err := NewClient[TestResource](ClientOptions{
			BaseURL: "http://crud-service/resource-path",
		})
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, baseURL, c.client.BaseURL.String())
	})

	t.Run("create new client with default headers to add in request", func(t *testing.T) {
		h := http.Header{}
		h.Set("Foo", "bar")
		h.Set("Taz", "ok")

		client, err := NewClient[TestResource](ClientOptions{
			BaseURL: baseURL,
			Headers: h,
		})
		require.NoError(t, err)
		require.NotNil(t, client)

		gock.New(baseURL).
			Get("export").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			BodyString("")

		_, err = client.Export(context.Background(), Options{})
		require.NoError(t, err)
	})
}

type NestedResource struct {
	Field string `json:"field"`
}

type TestResource struct {
	Field    string         `json:"field"`
	ID       string         `json:"_id"`
	IntField int            `json:"intField"`
	Nested   NestedResource `json:"nested"`
}

func TestExport(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	t.Run("export data", func(t *testing.T) {
		responseBody := `
		{"field": "v-1","intField":1,"_id":"my-id-1"}
		{"field": "v-2","intField":2,"_id":"my-id-2"}
		{"field": "v-3","intField":3,"_id":"my-id-3"}
		`

		gock.New(baseURL).
			Get("export").
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{})
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

		gock.New(baseURL).
			Get("export").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{Filter: filter})
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

		gock.New(baseURL).
			Get("export").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{Filter: filter})
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
		gock.New(baseURL).
			Get("export").
			Reply(500).
			AddHeader("Content-Type", "application/json").
			BodyString(`{"message":"error message"}`)

		resources, err := client.Export(ctx, Options{})
		require.EqualError(t, err, "error message")
		require.Nil(t, resources)
	})

	t.Run("export data and proxy headers in request", func(t *testing.T) {
		responseBody := `
		{"field": "v-1","intField":1,"_id":"my-id-1"}
		{"field": "v-2","intField":2,"_id":"my-id-2"}
		{"field": "v-3","intField":3,"_id":"my-id-3"}
		`

		gock.New(baseURL).
			Get("export").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			BodyString(responseBody)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.Export(ctx, Options{
			Headers: h,
		})
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
}

func TestGetById(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	id := "my-id-1"
	expectedElement := TestResource{
		Field:    "v-1",
		IntField: 1,
		ID:       id,
	}

	t.Run("get element by id", func(t *testing.T) {
		gock.New(baseURL).
			Get(id).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.GetByID(ctx, id, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("get element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
		}

		gock.New(baseURL).
			Get(id).
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.GetByID(ctx, id, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("throws - not found", func(t *testing.T) {
		gock.New(baseURL).
			Get(id).
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		resource, err := client.GetByID(ctx, id, Options{})
		require.EqualError(t, err, "element not found")
		require.Nil(t, resource)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		gock.New(baseURL).
			Get(id).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedElement)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.GetByID(ctx, id, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resources)
	})
}

func TestPatch(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	id := "my-id-1"
	expectedElement := TestResource{
		Field:    "v-1",
		IntField: 1,
		ID:       id,
		Nested: NestedResource{
			Field: "something",
		},
	}
	body := PatchBody{
		Set: map[string]any{
			"field":        "v-1",
			"nested.field": "something",
		},
	}
	expectedBody := `{"$set":{"field":"v-1","nested.field":"something"}}`

	t.Run("patch element", func(t *testing.T) {
		gock.New(baseURL).
			Patch(id).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchById(ctx, id, body, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("patch element with addToSet", func(t *testing.T) {
		gock.New(baseURL).
			Patch(id).
			BodyString(`{"$addToSet":{"something":{"$each":["a","b"]}}}`).
			Reply(200).
			JSON(expectedElement)

		type patchBodyAddSomething struct {
			Something any `json:"something"`
		}

		type eachOperatorBody struct {
			Each []string `json:"$each"`
		}

		body := PatchBody{
			AddToSet: patchBodyAddSomething{
				Something: eachOperatorBody{
					Each: []string{"a", "b"},
				},
			},
		}

		resource, err := client.PatchById(ctx, id, body, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("patch element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
		}

		gock.New(baseURL).
			Patch(id).
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchById(ctx, id, body, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("throws - not found", func(t *testing.T) {
		gock.New(baseURL).
			Patch(id).
			BodyString(expectedBody).
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		resource, err := client.PatchById(ctx, id, body, Options{})
		require.EqualError(t, err, "element not found")
		require.Nil(t, resource)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		gock.New(baseURL).
			Patch(id).
			BodyString(expectedBody).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedElement)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.PatchById(ctx, id, body, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resources)
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	id := "my-id-1"
	id2 := "my-id-2"
	expectedElements := []TestResource{
		{
			Field:    "v-1",
			IntField: 1,
			ID:       id,
			Nested: NestedResource{
				Field: "something",
			},
		},
		{
			Field:    "v-2",
			IntField: 2,
			ID:       id2,
		},
	}

	t.Run("list element", func(t *testing.T) {
		gock.New(baseURL).
			Get("").
			Reply(200).
			JSON(expectedElements)

		resource, err := client.List(ctx, id, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resource)
	})

	t.Run("list element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
			Skip:       4,
			Limit:      5,
		}

		gock.New(baseURL).
			Get("").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			JSON(expectedElements)

		resources, err := client.List(ctx, id, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resources)
	})

	t.Run("throws - not found", func(t *testing.T) {
		gock.New(baseURL).
			Get("").
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		resource, err := client.List(ctx, id, Options{})
		require.EqualError(t, err, "element not found")
		require.Nil(t, resource)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		gock.New(baseURL).
			Get("").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedElements)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.List(ctx, id, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resources)
	})
}

func getClient(t *testing.T) Client[TestResource] {
	t.Helper()

	client, err := NewClient[TestResource](ClientOptions{
		BaseURL: baseURL,
	})
	require.NoError(t, err)

	return client
}
