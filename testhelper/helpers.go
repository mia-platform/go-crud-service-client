package testhelper

import (
	"encoding/json"
	"testing"

	"github.com/mia-platform/go-crud-service-client/internal/types"
	"github.com/stretchr/testify/require"
)

type Filter types.Filter

func ParseResponseToNdjson[TResource any](t *testing.T, response []TResource) string {
	t.Helper()

	var responseBytes = []byte("")
	separator := []byte("\n")

	for _, elem := range response {
		marshElem, err := json.Marshal(elem)
		require.NoError(t, err, "Unexpected error to pars elem in NdJSON")

		responseBytes = append(responseBytes, separator...)
		responseBytes = append(responseBytes, marshElem...)
	}

	return string(responseBytes)
}
