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
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/mia-platform/go-crud-service-client/testhelper"

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

		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
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
		testhelper.NewGockScope(t, baseURL, http.MethodGet, id).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.GetByID(ctx, id, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("get element with filter", func(t *testing.T) {
		filter := Filter{
			Fields: map[string]string{
				"mockField": "mockValue",
			},
			Projection: []string{"field"},
		}

		testhelper.NewGockScope(t, baseURL, http.MethodGet, id).
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.GetByID(ctx, id, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, id).
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
		testhelper.NewGockScope(t, baseURL, http.MethodGet, id).
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
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "").
			Reply(200).
			JSON(expectedElements)

		resource, err := client.List(ctx, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resource)
	})

	t.Run("list element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
			Skip:       4,
			Limit:      5,
			Sort:       "field",
		}

		testhelper.NewGockScope(t, baseURL, http.MethodGet, "").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			JSON(expectedElements)

		resources, err := client.List(ctx, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resources)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "").
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		resource, err := client.List(ctx, Options{})
		require.EqualError(t, err, "element not found")
		require.Nil(t, resource)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedElements)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.List(ctx, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedElements, resources)
	})
}

func TestCount(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	expectedResult := 42

	t.Run("count elements", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "count").
			Reply(200).
			JSON(strconv.Itoa(expectedResult))

		actualResult, err := client.Count(ctx, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedResult, actualResult)
	})

	t.Run("count elements with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
			Skip:       4,
			Limit:      5,
		}

		testhelper.NewGockScope(t, baseURL, http.MethodGet, "count").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			JSON(strconv.Itoa(expectedResult))

		actualResult, err := client.Count(ctx, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, expectedResult, actualResult)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "count").
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		actualResult, err := client.Count(ctx, Options{})
		require.EqualError(t, err, "element not found")
		require.Equal(t, 0, actualResult)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "count").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(strconv.Itoa(expectedResult))

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		actualResult, err := client.Count(ctx, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedResult, actualResult)
	})
}

func TestExport(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	response := []TestResource{
		{Field: "v-1", IntField: 1, ID: "my-id-1"},
		{Field: "v-2", IntField: 2, ID: "my-id-2"},
		{Field: "v-3", IntField: 3, ID: "my-id-3"},
	}

	responseBody := testhelper.ParseResponseToNdjson[TestResource](t, response)

	t.Run("export data", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
			Reply(200).
			AddHeader("Content-Type", "application/x-ndjson").
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{})
		require.NoError(t, err)
		require.Equal(t, response, resources)
	})

	t.Run("export data with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
			MongoQuery: map[string]any{
				"field": map[string]any{
					"$in": []string{"v-1", "v-2"},
				},
			},
			Limit: 5,
		}

		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			AddHeader("Content-Type", "application/x-ndjson").
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, response, resources)
	})

	t.Run("export data with only mongo query", func(t *testing.T) {
		filter := Filter{
			MongoQuery: map[string]any{
				"field": map[string]any{"$in": []string{"v-1", "v-2"}},
			},
		}

		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			Reply(200).
			AddHeader("Content-Type", "application/x-ndjson").
			BodyString(responseBody)

		resources, err := client.Export(ctx, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, response, resources)
	})

	t.Run("throws with errors", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
			Reply(500).
			AddHeader("Content-Type", "application/json").
			BodyString(`{"message":"error message"}`)

		resources, err := client.Export(ctx, Options{})
		require.EqualError(t, err, "error message")
		require.Nil(t, resources)
	})

	t.Run("export data and proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodGet, "export").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			AddHeader("Content-Type", "application/x-ndjson").
			BodyString(responseBody)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		resources, err := client.Export(ctx, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, response, resources)
	})
}

func TestPatchById(t *testing.T) {
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
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, id).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchById(ctx, id, body, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("patch element with addToSet", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, id).
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

		testhelper.NewGockScope(t, baseURL, http.MethodPatch, id).
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchById(ctx, id, body, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, id).
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
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, id).
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

func TestPatchMany(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	expectedElement := TestResource{
		Field:    "v-1",
		IntField: 1,
		ID:       "id",
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
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "").
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchMany(ctx, body, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("patch element with addToSet", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "").
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

		resource, err := client.PatchMany(ctx, body, Options{})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("patch element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
		}

		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedElement)

		resource, err := client.PatchMany(ctx, body, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resource)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "").
			BodyString(expectedBody).
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "element not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		resource, err := client.PatchMany(ctx, body, Options{})
		require.EqualError(t, err, "element not found")
		require.Nil(t, resource)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "").
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

		resources, err := client.PatchMany(ctx, body, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, &expectedElement, resources)
	})
}

func TestPatchBulk(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	body := []PatchBulkItem{
		{
			Filter: PatchBulkFilter{
				Fields: map[string]string{
					"field": "v-1",
				},
			},
			Update: PatchBody{
				Set: map[string]any{
					"field":        "v-1",
					"nested.field": "something",
				},
			},
		},
		{
			Filter: PatchBulkFilter{
				Fields: map[string]string{
					"field": "v-2",
				},
			},
			Update: PatchBody{
				Set: map[string]any{
					"field":        "v-2",
					"nested.field": "another",
				},
			},
		},
	}
	expectedBody := `[{"filter":{"field":"v-1"},"update":{"$set":{"field":"v-1","nested.field":"something"}}},{"filter":{"field":"v-2"},"update":{"$set":{"field":"v-2","nested.field":"another"}}}]`
	expectedResponse := 3

	t.Run("patch element", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "bulk").
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedResponse)

		n, err := client.PatchBulk(ctx, body, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedResponse, n)
	})

	t.Run("patch element with addToSet", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "bulk").
			BodyString(`[{"filter":{"_q":"{\"foo\":\"bar\"}"},"update":{"$addToSet":{"something":{"$each":["a","b"]}}}}]`).
			Reply(200).
			JSON(3)

		type patchBodyAddSomething struct {
			Something any `json:"something"`
		}

		type eachOperatorBody struct {
			Each []string `json:"$each"`
		}

		body := []PatchBulkItem{
			{
				Filter: PatchBulkFilter{MongoQuery: map[string]any{"foo": "bar"}},
				Update: PatchBody{
					AddToSet: patchBodyAddSomething{
						Something: eachOperatorBody{
							Each: []string{"a", "b"},
						},
					},
				},
			},
		}

		n, err := client.PatchBulk(ctx, body, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedResponse, n)
	})

	t.Run("patch element with filter", func(t *testing.T) {
		filter := Filter{
			Projection: []string{"field"},
		}

		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "bulk").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			BodyString(expectedBody).
			Reply(200).
			JSON(expectedResponse)

		resource, err := client.PatchBulk(ctx, body, Options{Filter: filter})
		require.NoError(t, err)
		require.Equal(t, expectedResponse, resource)
	})

	t.Run("returns 0 if not found element", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "bulk").
			BodyString(expectedBody).
			Reply(200).
			JSON(0)

		n, err := client.PatchBulk(ctx, body, Options{})
		require.NoError(t, err)
		require.Equal(t, 0, n)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPatch, "bulk").
			BodyString(expectedBody).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedResponse)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		n, err := client.PatchBulk(ctx, body, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedResponse, n)
	})
}

func TestPatchBulkFilter(t *testing.T) {
	t.Run("marshal and unmarshal correctly", func(t *testing.T) {
		filter := PatchBulkFilter{
			Fields: map[string]string{
				"f1": "v1",
				"f2": "v2",
			},
			MongoQuery: map[string]any{
				"field": map[string]any{
					"$in": []any{"v-1", "v-2"},
				},
			},
		}

		f, err := json.Marshal(filter)
		require.NoError(t, err)
		require.JSONEq(t, `{"f1":"v1","f2":"v2","_q":"{\"field\":{\"$in\":[\"v-1\",\"v-2\"]}}"}`, string(f))

		actual := PatchBulkFilter{}
		err = json.Unmarshal(f, &actual)
		require.NoError(t, err)
		require.Equal(t, actual, filter)
	})

	t.Run("without mongo query", func(t *testing.T) {
		filter := PatchBulkFilter{
			Fields: map[string]string{
				"f1": "v1",
				"f2": "v2",
			},
		}

		f, err := json.Marshal(filter)
		require.NoError(t, err)
		require.JSONEq(t, `{"f1":"v1","f2":"v2"}`, string(f))

		actual := PatchBulkFilter{}
		err = json.Unmarshal(f, &actual)
		require.NoError(t, err)
		require.Equal(t, actual, filter)
	})

	t.Run("without fields", func(t *testing.T) {
		filter := PatchBulkFilter{
			MongoQuery: map[string]any{
				"field": map[string]any{
					"$in": []any{"v-1", "v-2"},
				},
			},
		}

		f, err := json.Marshal(filter)
		require.NoError(t, err)
		require.JSONEq(t, `{"_q":"{\"field\":{\"$in\":[\"v-1\",\"v-2\"]}}"}`, string(f))

		actual := PatchBulkFilter{}
		err = json.Unmarshal(f, &actual)
		require.NoError(t, err)
		require.Equal(t, actual, filter)
	})

	t.Run("empty filter", func(t *testing.T) {
		filter := PatchBulkFilter{}

		f, err := json.Marshal(filter)
		require.NoError(t, err)
		require.JSONEq(t, `{}`, string(f))

		actual := PatchBulkFilter{}
		err = json.Unmarshal(f, &actual)
		require.NoError(t, err)
		require.Equal(t, actual, filter)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	expectedID := "my _id"

	resourceToCreate := TestResource{
		Field:    "v-1",
		IntField: 1,
		Nested: NestedResource{
			Field: "something",
		},
	}

	expectedBody, err := json.Marshal(resourceToCreate)
	require.NoError(t, err)

	t.Run("creates resource", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "").
			BodyString(string(expectedBody)).
			Reply(200).
			JSON(map[string]string{
				"_id": expectedID,
			})

		actualID, err := client.Create(ctx, resourceToCreate, Options{})
		require.NoError(t, err)
		require.Equal(t, expectedID, actualID)
	})

	t.Run("throws - bad request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "").
			BodyString(string(expectedBody)).
			Reply(400).
			JSON(CrudErrorResponse{
				Message:    "missing required field",
				StatusCode: 400,
				Error:      "Bad Request",
			})

		actualID, err := client.Create(ctx, resourceToCreate, Options{})
		require.EqualError(t, err, "missing required field")
		require.Empty(t, actualID)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "").
			BodyString(string(expectedBody)).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(map[string]string{
				"_id": expectedID,
			})

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		actualID, err := client.Create(ctx, resourceToCreate, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedID, actualID)
	})
}

func TestCreateMany(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	expectedIDs := []map[string]string{
		{"_id": "1"},
		{"_id": "2"},
	}

	resourcesToCreate := []TestResource{
		{
			Field:    "v-1",
			IntField: 1,
			Nested: NestedResource{
				Field: "something",
			},
		},
		{
			Field:    "v-2",
			IntField: 2,
			Nested: NestedResource{
				Field: "something2",
			},
		},
	}

	expectedBody, err := json.Marshal(resourcesToCreate)
	require.NoError(t, err)

	t.Run("creates resource", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "bulk").
			BodyString(string(expectedBody)).
			Reply(200).
			JSON(expectedIDs)

		actualIDs, err := client.CreateMany(ctx, resourcesToCreate, Options{})
		require.NoError(t, err)
		require.Equal(t, []CreatedResource{
			{ID: "1"},
			{ID: "2"},
		}, actualIDs)
	})

	t.Run("throws - bad request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "bulk").
			BodyString(string(expectedBody)).
			Reply(400).
			JSON(CrudErrorResponse{
				Message:    "missing required field",
				StatusCode: 400,
				Error:      "Bad Request",
			})

		actualID, err := client.CreateMany(ctx, resourcesToCreate, Options{})
		require.EqualError(t, err, "missing required field")
		require.Empty(t, actualID)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "bulk").
			BodyString(string(expectedBody)).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).
			JSON(expectedIDs)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		actualIDs, err := client.CreateMany(ctx, resourcesToCreate, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, []CreatedResource{
			{ID: "1"},
			{ID: "2"},
		}, actualIDs)
	})
}

func TestDeleteById(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	id := "my-id-1"

	t.Run("delete element", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, id).
			Reply(204)

		err := client.DeleteById(ctx, id, Options{})
		require.NoError(t, err)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, id).
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		err := client.DeleteById(ctx, id, Options{})
		require.EqualError(t, err, "not found")
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, id).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(204)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		err := client.DeleteById(ctx, id, Options{
			Headers: h,
		})
		require.NoError(t, err)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	t.Run("delete element", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, "").
			Reply(200).BodyString("3")

		n, err := client.DeleteMany(ctx, Options{})
		require.NoError(t, err)
		require.Equal(t, 3, n)
	})

	t.Run("throws - not found", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, "").
			Reply(404).
			JSON(CrudErrorResponse{
				Message:    "not found",
				StatusCode: 404,
				Error:      "Not Found",
			})

		n, err := client.DeleteMany(ctx, Options{})
		require.EqualError(t, err, "not found")
		require.Equal(t, 0, n)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodDelete, "").
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			Reply(200).BodyString("4")

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		n, err := client.DeleteMany(ctx, Options{
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, 4, n)
	})
}

func TestUpsertOne(t *testing.T) {
	ctx := context.Background()
	client := getClient(t)

	expectedElement := &TestResource{
		Field:    "v-1",
		IntField: 1,
		ID:       "my-id",
		Nested: NestedResource{
			Field: "something",
		},
	}
	filter := Filter{
		MongoQuery: map[string]any{
			"field": "v-1",
		},
	}
	requestBody := UpsertBody{
		SetOnInsert: map[string]any{
			"field": "v-1",
			"id":    "my-id",
		},
		Set: map[string]any{
			"intField": 1,
			"nested": map[string]any{
				"field": "something",
			},
		},
	}

	t.Run("upsert one element", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "upsert-one").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			JSON(requestBody).
			Reply(200).
			JSON(expectedElement)

		response, err := client.UpsertOne(ctx, requestBody, Options{
			Filter: filter,
		})
		require.NoError(t, err)
		require.Equal(t, expectedElement, response)
	})

	t.Run("throws if crud fails", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "upsert-one").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			JSON(requestBody).
			Reply(500).
			JSON(CrudErrorResponse{
				Message:    "broken",
				StatusCode: 500,
				Error:      "Internal Server Error",
			})

		response, err := client.UpsertOne(ctx, requestBody, Options{
			Filter: filter,
		})
		require.EqualError(t, err, "broken")
		require.Nil(t, response)
	})

	t.Run("proxy headers in request", func(t *testing.T) {
		testhelper.NewGockScope(t, baseURL, http.MethodPost, "upsert-one").
			AddMatcher(testhelper.CrudQueryMatcher(t, testhelper.Filter(filter))).
			MatchHeaders(map[string]string{
				"foo": "bar",
				"taz": "ok",
			}).
			JSON(requestBody).
			Reply(200).
			JSON(expectedElement)

		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("taz", "ok")

		response, err := client.UpsertOne(ctx, requestBody, Options{
			Filter:  filter,
			Headers: h,
		})
		require.NoError(t, err)
		require.Equal(t, expectedElement, response)
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
