package test

import (
	"github.com/simon987/ws_bucket/api"
	"testing"
)

func TestCreateClient(t *testing.T) {

	r := createClient(api.CreateClientRequest{
		Alias: "testcreateclient",
	})

	if r.Ok != true {
		t.Error()
	}
}

func createClient(request api.CreateClientRequest) (ar *api.CreateClientResponse) {

	resp := Post("/client", request)
	UnmarshalResponse(resp, &ar)
	return
}
