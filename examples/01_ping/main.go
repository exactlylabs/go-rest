package main

import (
	"github.com/exactlylabs/go-rest/pkg/restapi"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
)

func main() {
	server, err := restapi.NewWebServer()
	if err != nil {
		panic(err)
	}
	server.Route("/ping", Ping, "GET")
	server.Run("127.0.0.1:5000")
}

func Ping(c *webcontext.Context) {
	c.JSON(200, map[string]string{"response": "pong"})
}
