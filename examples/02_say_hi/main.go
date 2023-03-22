/*
	This example shows how to grab the Query Parameters, parse a JSON from the Body and also how to add field validation errors
*/
package main

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/exactlylabs/go-errors/pkg/errors"

	"github.com/exactlylabs/go-rest/pkg/restapi"
	"github.com/exactlylabs/go-rest/pkg/restapi/apierrors"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
	"github.com/go-playground/validator/v10"
)

func main() {
	server, err := restapi.NewWebServer()
	if err != nil {
		panic(err)
	}
	server.Route("/", Hi, "GET")
	server.Route("/", HiPost, "POST")
	server.Run("127.0.0.1:5000")
}

func Hi(ctx *webcontext.Context) {
	if !ctx.QueryParams().Has("name") {
		// This will register an error that is going to be returned to the caller
		// Missing field is a common error, therefore we have a generic error for that
		ctx.AddFieldError("name", apierrors.MissingFieldError)
		return
	} else if ctx.QueryParams().Get("name") == "" {
		ctx.AddFieldError("name", apierrors.SingleFieldError("cannot be empty", "empty"))
		return
	}
	name := ctx.QueryParams().Get("name")
	ctx.JSON(200, map[string]string{"response": "Hi " + name})
}

type HiPostArgs struct {
	Name string `validate:"required" json:"name"`
}

func HiPost(ctx *webcontext.Context) {
	// For POST, it's possible to use json.Decode(ctx.Request.Body), or a validator such as go-playground/validator
	// In both cases, field errors are a bit harder right now, but we may include our own validator in the future
	args := &HiPostArgs{}
	json.NewDecoder(ctx.Request.Body).Decode(args)
	// This could be a global variable, so you just do it once, or a dependency (see dependency example)
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := validate.Struct(args)
	ve := validator.ValidationErrors{}
	if ok := errors.As(err, &ve); ok {
		for _, err := range ve {
			ctx.AddFieldError(err.Field(), apierrors.SingleFieldError(err.Error(), err.ActualTag()))
		}
		return
	}
	ctx.JSON(200, map[string]string{"response": "Hi " + args.Name})
}
