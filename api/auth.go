package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"math/rand"
)

func (api *WebApi) CreateClient(ctx *fasthttp.RequestCtx) {

	//TODO: auth

	req := &CreateClientRequest{}
	err := json.Unmarshal(ctx.Request.Body(), req)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		Json(CreateClientResponse{
			Ok: false,
		}, ctx)
		return
	}

	if !req.IsValid() {
		ctx.Response.SetStatusCode(400)
		Json(CreateClientResponse{
			Ok: false,
		}, ctx)
		return
	}

	client := api.createClient(req)

	Json(CreateClientResponse{
		Ok:     true,
		Secret: client.Secret,
	}, ctx)
}

func (api *WebApi) createClient(req *CreateClientRequest) *Client {

	client := &Client{
		Alias:  req.Alias,
		Secret: genSecret(),
	}

	api.db.Create(client)

	logrus.WithFields(logrus.Fields{
		"client": client,
	}).Info("Created client")

	return client
}

func genSecret() string {
	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		bytes[i] = byte(48 + rand.Intn(122-48))
	}
	return string(bytes)
}
