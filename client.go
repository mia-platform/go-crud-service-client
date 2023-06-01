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
	"github.com/mia-platform/go-crud-service-client/internal/types"
)

type Client[Resource any] struct {
	client *jsonclient.Client
}

type ClientOptions struct {
	BaseURL string
}

type Filter types.Filter

func NewClient[Resource any](options ClientOptions) (Client[Resource], error) {
	client, err := jsonclient.New(jsonclient.Options{
		BaseURL: options.BaseURL,
	})
	if err != nil {
		return Client[Resource]{}, fmt.Errorf("%w: %s", ErrCreateClient, err)
	}
	return Client[Resource]{
		client: client,
	}, err
}

func (c Client[Resource]) Export(ctx context.Context, path string, filter Filter) ([]Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodGet, "export", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateRequest, err)
	}

	if err = addCrudQuery(req, types.Filter(filter)); err != nil {
		return nil, err
	}

	responseBuf := bytes.NewBuffer(nil)
	_, err = c.client.Do(req, responseBuf)
	if err != nil {
		return nil, err
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
