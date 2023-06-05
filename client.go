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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/davidebianchi/go-jsonclient"
)

type Client[Resource any] struct {
	client *jsonclient.Client
}

// NewClient create a new client to interact with crud-service
func NewClient[Resource any](options ClientOptions) (Client[Resource], error) {
	client, err := jsonclient.New(jsonclient.Options{
		BaseURL: options.BaseURL,
		Headers: options.convertHeaders(),
	})
	if err != nil {
		return Client[Resource]{}, fmt.Errorf("%w: %s", ErrCreateClient, err)
	}
	return Client[Resource]{
		client: client,
	}, err
}

// Export calls /export endpoint of crud-service. It is possible to add filters.
// Exports does not have max limits.
func (c Client[Resource]) Export(ctx context.Context, options Options) ([]Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodGet, "export", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateRequest, err)
	}

	if err := options.setOptionsInRequest(req); err != nil {
		return nil, err
	}

	responseBuf := bytes.NewBuffer(nil)
	_, err = c.client.Do(req, responseBuf)
	if err != nil {
		return nil, responseError(err)
	}

	resources := []Resource{}

	decoder := json.NewDecoder(responseBuf)
	for decoder.More() {
		resource := new(Resource)
		if err := decoder.Decode(resource); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		resources = append(resources, *resource)
	}

	return resources, nil
}

// GetById get a resource by _id
func (c Client[Resource]) GetByID(ctx context.Context, id string, options Options) (*Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodGet, id, nil)
	if err != nil {
		return nil, err
	}

	if err := options.setOptionsInRequest(req); err != nil {
		return nil, err
	}

	resource := new(Resource)
	if _, err := c.client.Do(req, resource); err != nil {
		return nil, responseError(err)
	}
	return resource, nil
}
