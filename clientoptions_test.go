package crud

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertHeaders(t *testing.T) {
	t.Run("convert from http.Headers to map", func(t *testing.T) {
		h := http.Header{}
		h.Set("foo", "bar")
		h.Set("key1", "val1")

		clientOpts := ClientOptions{
			Headers: h,
		}

		res := clientOpts.convertHeaders()

		require.Equal(t, map[string]string{
			"Foo":  "bar",
			"Key1": "val1",
		}, res)
	})
}
