package mocks

import "net/http"

type MockClient struct{}

var (
	DoFunc func(req *http.Request) (*http.Response, error)
)

//Do implements the HTTPClient interface in clients/restclient/restclient.go
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
