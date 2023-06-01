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

func addCrudQuery(req *http.Request, filter types.Filter) error {
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
