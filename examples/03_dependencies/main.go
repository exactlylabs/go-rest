/*
	This example shows the usage of Dependencies
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

// We can register Dependency providers in the server
// These providers will be available for all route handlers and are great for testing purposes
// In this case, instead of using a global validate, we use one provided by this function
// it could just return a global one, without a problem, but by using injection, we ensure the handler has everything it needs
func ValidatorProvider(ctx *webcontext.Context) any {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

func main() {
	server, err := restapi.NewWebServer()
	if err != nil {
		panic(err)
	}
	server.Route("/", HiPost, "POST")
	// We register a dependency by passing its factory and an object that this factory returns
	// Note: Generics won't work here, so we need this for reflection
	server.AddDependency(ValidatorProvider, &validator.Validate{})
	server.Run("127.0.0.1:5000")
}

type HiPostArgs struct {
	Name string `validate:"required" json:"name"`
}

// See the addition of a new argument here. With this argument, we are telling the server that this handler needs a *validator.Validate dependency injected in it
// If there's a dependency provider registered that matches this variable type, then it is going to be injected here
func HiPost(ctx *webcontext.Context, validate *validator.Validate) {
	// For POST, it's possible to use json.Decode(ctx.Request.Body), or a validator such as go-playground/validator
	// In both cases, field errors are a bit harder right now, but we may include our own validator in the future
	args := &HiPostArgs{}
	json.NewDecoder(ctx.Request.Body).Decode(args)
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
