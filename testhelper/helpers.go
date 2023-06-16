package testhelper

import (
	"encoding/json"

	"github.com/mia-platform/go-crud-service-client/internal/types"
)

type Filter types.Filter

func ParseResponseToNdjson[TResource any](response []TResource) (string, error) {
	var responseBytes = []byte("")
	separator := []byte("\n")

	for _, elem := range response {
		marshElem, err := json.Marshal(elem)
		if err != nil {
			return "", err
		}

		responseBytes = append(responseBytes, separator...)
		responseBytes = append(responseBytes, marshElem...)
	}

	return string(responseBytes), nil
}
