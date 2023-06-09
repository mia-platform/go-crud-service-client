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
	"fmt"
	"net/http"
	"testing"

	"github.com/davidebianchi/go-jsonclient"
	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	t.Run("with invalid json http error response", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Raw = []byte(`{"invalid json"}`)

		err := responseError(httpErr)

		require.EqualError(t, err, "invalid character '}' after object key")
	})

	t.Run("correctly wrap HTTPErrors", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		err := responseError(httpErr)

		require.ErrorIs(t, err, ErrResponse)
		require.EqualError(t, err, "Some message")

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:   httpErr.Response,
			StatusCode: httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{
				Message:    "Some message",
				Error:      "my error",
				StatusCode: 500,
			},
			Err: ErrResponse,
			Raw: httpErr.Raw,
		}, crudError)
	})

	t.Run("with invalid json http error response - content-type with charset", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Response.Header.Set("Content-Type", "application/json; charset=utf-8")

		err := responseError(httpErr)

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:   httpErr.Response,
			StatusCode: httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{
				Message:    "Some message",
				Error:      "my error",
				StatusCode: 500,
			},
			Err: ErrResponse,
			Raw: httpErr.Raw,
		}, crudError)
		require.EqualError(t, err, "Some message")
	})

	t.Run("with content-type set to text/html", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Response.Header.Set("Content-Type", "text/html")
		httpErr.Raw = []byte(`a raw message`)

		err := responseError(httpErr)

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:     httpErr.Response,
			StatusCode:   httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{},
			Err:          ErrResponse,
			Raw:          httpErr.Raw,
		}, crudError)
		require.EqualError(t, err, "a raw message")
	})

	t.Run("without content-type header", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Response.Header.Del("Content-Type")
		httpErr.Raw = []byte(`a raw message`)

		err := responseError(httpErr)

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:     httpErr.Response,
			StatusCode:   httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{},
			Err:          ErrResponse,
			Raw:          httpErr.Raw,
		}, crudError)
		require.EqualError(t, err, "a raw message")
	})

	t.Run("without body - content-type json", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Raw = []byte(``)

		err := responseError(httpErr)

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:     httpErr.Response,
			StatusCode:   httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{},
			Err:          ErrResponse,
			Raw:          httpErr.Raw,
		}, crudError)
		require.EqualError(t, err, "error body from crud-service is empty")
	})

	t.Run("without body - content-type not json", func(t *testing.T) {
		httpErr := getJsonClientHttpError()
		httpErr.Response.Header.Del("Content-Type")
		httpErr.Raw = []byte(``)

		err := responseError(httpErr)

		crudError := &HTTPError{}
		require.ErrorAs(t, err, &crudError)
		require.Equal(t, &HTTPError{
			Response:     httpErr.Response,
			StatusCode:   httpErr.StatusCode,
			ResponseBody: CrudErrorResponse{},
			Err:          ErrResponse,
			Raw:          httpErr.Raw,
		}, crudError)
		require.EqualError(t, err, "error body from crud-service is empty")
	})

	t.Run("response error is not a jsonclient.HTTPError", func(t *testing.T) {
		someErr := fmt.Errorf("an error message")
		err := responseError(someErr)

		require.EqualError(t, err, "an error message")
	})
}

func getJsonClientHttpError() *jsonclient.HTTPError {
	response := &http.Response{
		Header: http.Header{},
	}
	response.Header.Set("Content-Type", "application/json")

	h := &jsonclient.HTTPError{
		Response:   response,
		StatusCode: 500,
		Err:        ErrResponse,
		Raw:        []byte(`{"message":"Some message","statusCode":500,"error":"my error"}`),
	}

	return h
}
