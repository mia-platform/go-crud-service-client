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

// Body to update a document in the collection
type PatchBody struct {
	// Set replaces the value of the field with specified value. It is possible also
	// to use with nested fields: e.g. `"a.b": "update"`
	Set any `json:"$set,omitempty"`
	// Unset a particular document value
	Unset map[string]bool `json:"$unset,omitempty"`
	// Inc increment a field by a specified value
	Inc map[string]int `json:"$inc,omitempty"`
	// Mul multiply the value of a field by a specified number
	Mul map[string]int `json:"$mul,omitempty"`
	// CurrentDate sets the value of a field to the current date. The field MUST
	// be of type Date
	CurrentDate any `json:"$currentDate,omitempty"`
	// Push appends a value to an array field
	Push any `json:"$push,omitempty"`
	// Pull removes a specified value from an array field
	Pull any `json:"$pull,omitempty"`
	// AddToSet appends a specified value to an array field unless the value is
	// already present
	AddToSet any `json:"$addToSet,omitempty"`
}

// PatchById update an element using commands in PatchBody
func (c Client[Resource]) PatchById(ctx context.Context, id string, body PatchBody, options Options) (*Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodPatch, id, body)
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

// List get the resources of the collection with the specified filter. It is limited by default
// and with a max page of 200 elements (by default).
// If you want to take more elements, use pagination
func (c Client[Resource]) List(ctx context.Context, options Options) ([]Resource, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}

	if err := options.setOptionsInRequest(req); err != nil {
		return nil, err
	}

	resources := []Resource{}
	if _, err := c.client.Do(req, &resources); err != nil {
		return nil, responseError(err)
	}
	return resources, nil
}

// The type that represents a newly created resource
type CreatedResource struct {
	ID string `json:"_id"`
}

// Create performs a POST request to create a new resource on the target crud. Returns the
// identifier of the created resource and any error that occurred.
func (c Client[Resource]) Create(ctx context.Context, resource Resource, options Options) (string, error) {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodPost, "", resource)
	if err != nil {
		return "", err
	}

	if err := options.setOptionsInRequest(req); err != nil {
		return "", err
	}

	var createdResource CreatedResource
	if _, err := c.client.Do(req, &createdResource); err != nil {
		return "", responseError(err)
	}
	return createdResource.ID, nil
}

// DeleteById deletes an element using the resource _id.
func (c Client[Resource]) DeleteById(ctx context.Context, id string, options Options) error {
	req, err := c.client.NewRequestWithContext(ctx, http.MethodDelete, id, nil)
	if err != nil {
		return err
	}

	if err := options.setOptionsInRequest(req); err != nil {
		return err
	}

	resource := new(Resource)
	if _, err := c.client.Do(req, resource); err != nil {
		return responseError(err)
	}
	return nil
}
