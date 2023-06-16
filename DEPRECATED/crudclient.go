/*
 * Copyright © 2020-present Mia s.r.l.
 * All rights reserved
 */

package crudclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"rbac-manager-bff/helpers"

	"github.com/davidebianchi/go-jsonclient"
)

// CRUD struct.
type CRUD struct {
	httpClient *jsonclient.Client
}

// New creates new CRUDClient.
func New(apiURL string) (*CRUD, error) {
	apiURLWithoutFinalSlash := strings.TrimSuffix(apiURL, "/")

	opts := jsonclient.Options{
		BaseURL: fmt.Sprintf("%s/", apiURLWithoutFinalSlash),
	}
	httpClient, err := jsonclient.New(opts)
	if err != nil {
		return nil, err
	}

	crudClient := &CRUD{
		httpClient,
	}
	return crudClient, nil
}

var ErrForbidden = errors.New("forbidden response")

func wrapError(err error) error {
	var httpError *jsonclient.HTTPError
	if errors.As(err, &httpError) && httpError.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: %s", ErrForbidden, err.Error())
	}
	return err
}

// Get fetch items based on a query from CRUD.
func (crud CRUD) Get(ctx context.Context, queryParam string, responseBody interface{}) error {
	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodGet, "?"+queryParam, nil)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBody); err != nil {
		return wrapError(err)
	}
	return nil
}

// Count items based on a query from CRUD.
func (crud CRUD) Count(ctx context.Context, queryParam string) (int, error) {
	var responseBuffer = &bytes.Buffer{}
	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodGet, "count?"+queryParam, nil)
	if err != nil {
		return 0, err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBuffer); err != nil {
		return 0, err
	}

	return strconv.Atoi(responseBuffer.String())
}

func (crud CRUD) Post(ctx context.Context, body interface{}, responseBody interface{}) error {
	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodPost, "", body)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	_, err = crud.httpClient.Do(req, responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (crud CRUD) Delete(ctx context.Context, queryParam string, body interface{}, responseBody interface{}) error {
	query := ""
	if queryParam != "" {
		query = "?" + queryParam
	}

	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodDelete, query, body)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBody); err != nil {
		return err
	}
	return nil
}

func (crud CRUD) DeleteById(ctx context.Context, id string, body interface{}, responseBody interface{}) error {

	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodDelete, id, body)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBody); err != nil {
		return err
	}
	return nil
}

func (crud CRUD) PatchBulkWithQuery(ctx context.Context, query string, requestBody interface{}, responseBody interface{}) error {
	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodPatch, query, requestBody)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBody); err != nil {
		return err
	}
	return nil
}

func (crud CRUD) PatchBulk(ctx context.Context, requestBody interface{}, responseBody interface{}) error {
	req, err := crud.httpClient.NewRequestWithContext(ctx, http.MethodPatch, "", requestBody)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, responseBody); err != nil {
		return err
	}
	return nil
}

// IsHealthy checks if crud is healthy.
func (crud CRUD) IsHealthy(ctx context.Context) error {
	req, err := crud.httpClient.NewRequest(http.MethodGet, "/-/healthz", nil)
	if err != nil {
		return err
	}

	helpers.SetHeadersToProxy(ctx, req.Header)

	if _, err := crud.httpClient.Do(req, nil); err != nil {
		return err
	}
	return nil
}
