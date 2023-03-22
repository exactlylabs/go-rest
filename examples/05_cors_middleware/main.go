/*
	This example shows you how to add CORS to your server, as well as any other middleware
*/
package main

import (
	"github.com/exactlylabs/go-rest/pkg/restapi"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
	"github.com/gorilla/handlers"
)

func main() {
	api, err := restapi.NewWebServer()
	if err != nil {
		panic(err)
	}
	api.AddMiddlewares(
		handlers.CORS(
			handlers.AllowedOrigins([]string{"http://127.0.0.1:3000"}), // use * for all origins
		),
	)
	api.Route("/ping", Ping)
	api.Run("127.0.0.1:5000")
}

func Ping(c *webcontext.Context) {
	c.JSON(200, map[string]string{"response": "pong"})
}
