package providers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type MockServerConfig struct {
	responseStatusCode int
	contentType        string
	headers            map[string]string
	response           string
}

func NewMockServerConfig(
	responseStatusCode int,
	contentType string,
	headers map[string]string,
	response string,
) *MockServerConfig {
	return &MockServerConfig{
		responseStatusCode: responseStatusCode,
		contentType:        contentType,
		headers:            headers,
		response:           response,
	}
}

func (config MockServerConfig) CreateMockServer() *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", config.contentType)
		w.WriteHeader(config.responseStatusCode)
		fmt.Fprintln(w, config.response)

		for header, value := range config.headers {
			w.Header().Set(header, value)
		}
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}
