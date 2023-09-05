package backend

import (
	"fmt"
	"net/http"
)

type Backend struct {
	Addr   string
	client http.Client
}

func New(addr string) *Backend {
	return &Backend{
		Addr:   addr,
		client: *http.DefaultClient,
	}
}

func (b *Backend) Handle(req *http.Request) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", b.Addr, req.RequestURI)
	backendRequest, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}
	return b.client.Do(backendRequest)
}

func (b *Backend) CheckHealth(url string) bool {
	res, err := http.DefaultClient.Get(fmt.Sprintf("%s%s", b.Addr, url))
	if err != nil {
		return false
	}
	if res.StatusCode != 200 {
		return false
	}
	return true
}
