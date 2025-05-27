//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/http_client_mock.go -package=mocks -source ./http_client.go

package domain

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
