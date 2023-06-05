package crud

import "net/http"

type ClientOptions struct {
	BaseURL string
	Headers http.Header
}

func (options ClientOptions) convertHeaders() map[string]string {
	h := map[string]string{}
	for name := range options.Headers {
		h[name] = options.Headers.Get(name)
	}
	return h
}
