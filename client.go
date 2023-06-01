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
	"net/url"
	"strconv"
	"strings"

	"github.com/davidebianchi/go-jsonclient"
)

type CrudClient[Resource any] struct {
	client *jsonclient.Client
}

type ClientOptions struct {
	BaseURL string
}

type Filter struct {
	MongoQuery map[string]any
	Limit      int
	Projection []string
}

func New[Resource any](options ClientOptions) (CrudClient[Resource], error) {
	client, err := jsonclient.New(jsonclient.Options{
		BaseURL: options.BaseURL,
	})
	if err != nil {
		return CrudClient[Resource]{}, fmt.Errorf("%w: %s", ErrCreateClient, err)
	}
	return CrudClient[Resource]{
		client: client,
	}, err
}

func (c CrudClient[Resource]) Export(ctx context.Context, path string, filter Filter) ([]Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodGet, "export", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCreateRequest, err)
	}

	if err = addCrudQuery(req, filter); err != nil {
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

func addCrudQuery(req *http.Request, filter Filter) error {
	query := url.Values{}
	if filter.MongoQuery != nil {
		queryBytes, err := json.Marshal(filter.MongoQuery)
		if err != nil {
			return err
		}
		query.Set("_q", string(queryBytes))
	}

	if filter.Limit != 0 {
		query.Set("_l", strconv.Itoa(filter.Limit))
	}

	if filter.Projection != nil {
		query.Set("_p", strings.Join(filter.Projection, ","))
	}

	req.URL.RawQuery = query.Encode()

	return nil
}
