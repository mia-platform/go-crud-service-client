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
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mia-platform/go-crud-service-client/internal/types"
)

type Filter types.Filter

type Options struct {
	Filter  Filter
	Headers http.Header
}

func (o Options) setOptionsInRequest(req *http.Request) error {
	if err := addCrudQueryToRequest(req, types.Filter(o.Filter)); err != nil {
		return err
	}

	addHeaderToRequest(req, o.Headers)

	return nil
}

func addCrudQueryToRequest(req *http.Request, filter types.Filter) error {
	query := url.Values{}
	if err := calculateFilter(query, filter); err != nil {
		return err
	}

	req.URL.RawQuery = query.Encode()

	return nil
}

func addHeaderToRequest(req *http.Request, headers http.Header) {
	for name := range headers {
		req.Header.Set(name, headers.Get(name))
	}
}

type Setter interface {
	Set(k, v string)
}

func calculateFilter(query Setter, filter types.Filter) error {
	if filter.Fields != nil {
		for field, value := range filter.Fields {
			query.Set(field, value)
		}
	}

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

	if filter.Skip != 0 {
		query.Set("_sk", strconv.Itoa(filter.Skip))
	}

	if filter.Sort != "" {
		query.Set("_s", filter.Sort)
	}

	return nil
}
