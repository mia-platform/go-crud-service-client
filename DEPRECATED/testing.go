/*
 * Copyright © 2022-present Mia s.r.l.
 * All rights reserved
 */

package crudclient

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

// MockUpsertWithQueryParameters mocks upsert to a collection.
func getHeadersMap(headers http.Header) map[string]string {
	requestHeadersMap := map[string]string{}
	if len(headers) != 0 {
		for name, values := range headers {
			requestHeadersMap[name] = values[0]
		}
	}
	return requestHeadersMap
}

// MockGetByID mocks get by id in a collection.
func MockGetByID(t *testing.T, baseURL string, statusCode int, id string, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Get(id).
		Reply(statusCode).
		JSON(responseBody)
}

// MockGet mocks get in a collection.
func MockGet(t *testing.T, baseURL string, statusCode int, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

func alwaysMatch(req *http.Request, greq *gock.Request) (bool, error) {
	return true, nil
}

// MockExport mocks get in a collection.
func MockExport(t *testing.T, baseURL string, statusCode int, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		Get("/export").
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

// Real mock that return a ndjson response
func MockExportWithNdjsonResponse[TResource any](t *testing.T, baseURL string, statusCode int, responseBody []TResource, headersToProxy http.Header, customMatcher gock.MatchFunc) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	var matcher gock.MatchFunc
	if customMatcher == nil {
		matcher = alwaysMatch
	} else {
		matcher = customMatcher
	}

	var responseBytes = []byte("")
	separator := []byte("\n")

	for _, elem := range responseBody {
		marshElem, err := json.Marshal(elem)
		require.Nil(t, err, "Unexpected error to marshal elem in MockExportWithResponseString")

		responseBytes = append(responseBytes, separator...)
		responseBytes = append(responseBytes, marshElem...)
	}

	gock.New(baseURL).
		Get("/export").
		AddMatcher(matcher).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		AddHeader("Content-Type", "application/x-ndjson").
		BodyString(string(responseBytes))
}

// MockPost mocks post in a collection.
func MockPost(t *testing.T, baseURL string, statusCode int, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

// MockDeleteByID mocks delete by id in a collection.
func MockDeleteByID(t *testing.T, baseURL string, statusCode int, id string, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Delete(id).
		Reply(statusCode).
		JSON(responseBody)
}

// MockPatchByIDWithBodyMatcher mocks patch by id in a collection with a custom body matcher.
func MockPatchByIDWithBodyMatcher(t *testing.T, baseURL string, statusCode int, id string, responseBody interface{}, headersToProxy http.Header, bodyMatcher gock.MatchFunc) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gockReq := gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Patch(id)

	if bodyMatcher != nil {
		gockReq.AddMatcher(bodyMatcher)
	}

	gockReq.
		Reply(statusCode).
		JSON(responseBody)
}

// MockPatchByID mocks patch by id in a collection.
func MockPatchByID(t *testing.T, baseURL string, statusCode int, id string, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	MockPatchByIDWithBodyMatcher(t, baseURL, statusCode, id, responseBody, headersToProxy, nil)
}

// MockDelete mocks delete in a collection.
func MockDelete(t *testing.T, baseURL string, statusCode int, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

// MockPatchBulk mocks post in a collection.
func MockPatchBulk(t *testing.T, baseURL string, statusCode int, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

func MockPatchBulkWithMatch(t *testing.T, baseURL string, statusCode int, matchQuery string, responseBody interface{}, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	requestUrl := strings.Join([]string{baseURL, matchQuery}, "?")
	gock.New(requestUrl).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Reply(statusCode).
		JSON(responseBody)
}

// MockIsHealthy mock the healthy function
func MockIsHealthy(t *testing.T, baseURL string, statusCode int, headersToProxy http.Header) {
	t.Helper()
	t.Cleanup(func() {
		gockCleanup(t)
	})
	gock.DisableNetworking()

	responseBody := map[string]interface{}{
		"status": "OK",
	}
	if statusCode >= 300 {
		responseBody = map[string]interface{}{
			"status": "KO",
		}
	}

	gock.New(baseURL).
		MatchHeaders(getHeadersMap(headersToProxy)).
		Get("/-/healthz").
		Reply(statusCode).
		JSON(responseBody)
}

func gockCleanup(t *testing.T) {
	t.Helper()

	if !gock.IsDone() {
		gock.OffAll()
		t.Fatal("fails to mock crud")
	}
	gock.Off()
}
