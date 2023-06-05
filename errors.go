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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/davidebianchi/go-jsonclient"
)

var (
	ErrCreateClient  = fmt.Errorf("fails to create client")
	ErrCreateRequest = fmt.Errorf("fails to create requests")

	ErrResponse = fmt.Errorf("crud error")
)

type HTTPError struct {
	Response     *http.Response
	StatusCode   int
	Err          error
	ResponseBody CrudErrorResponse
	Raw          []byte
}

func (e *HTTPError) Error() string {
	message := e.ResponseBody.Message
	if message == "" {
		message = string(e.Raw)
	}
	if message == "" {
		message = "error body from crud-service is empty"
	}

	return message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

type CrudErrorResponse struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	Error      string `json:"error,omitempty"`
}

func responseError(resErr error) error {
	var httpError *jsonclient.HTTPError

	if !errors.As(resErr, &httpError) {
		return resErr
	}

	errorResponse := CrudErrorResponse{}
	if strings.HasPrefix(httpError.Response.Header.Get("Content-Type"), "application/json") && string(httpError.Raw) != "" {
		if err := httpError.Unmarshal(&errorResponse); err != nil {
			return err
		}
	}

	return &HTTPError{
		Response:     httpError.Response,
		StatusCode:   httpError.StatusCode,
		Err:          ErrResponse,
		ResponseBody: errorResponse,
		Raw:          httpError.Raw,
	}
}
