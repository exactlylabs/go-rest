/*
This example shows you how to configure Sentry in your server

If you have Sentry, replace your dsn and run the example. You should see the error in your Sentry dashboard.

The example provides 3 endpoints. 1 panicking directly in the handler, 1 using the default errors package and another using exactlylabs/go-errors package
*/
package main

import (
	"flag"

	default_errors "github.com/pkg/errors"

	"github.com/exactlylabs/go-errors/pkg/errors"
	"github.com/exactlylabs/go-monitor/pkg/sentry"
	"github.com/exactlylabs/go-rest/pkg/restapi"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
)

// ErrBase is the base error for this package, it's useful if you want to differentiate errors from this package and third party packages
var ErrBase = errors.NewWithType("My Package Error", "MyPkgError")

// ErrInvalid is an example of a sentinel error that belongs to ErrBase
var ErrInvalid = errors.WrapWithType(ErrBase, "something went wrong in the API", "APIError")

func main() {
	dsn := flag.String("sentry-dsn", "", "Sentry DSN")
	flag.Parse()
	// To add Sentry, you just need to setup it and then replace the default recovery middleware by exactlylabs sentry middleware
	sentry.Setup(*dsn, "0.0.1", "testing", "Example06")
	// Always remember to add this when entering any goroutine, otherwise, the panic will not be captured
	defer sentry.NotifyIfPanic()

	api, err := restapi.NewWebServer(
		restapi.WithRecoveryMiddleware(sentry.SentryMiddleware),
	)
	if err != nil {
		panic(err)
	}

	api.Route("", DirectPanic)
	api.Route("/without-go-errors", WithoutGoErrors)
	api.Route("/with-go-errors", WithGoErrors)
	api.Run("127.0.0.1:5000")
}

func DirectPanic(c *webcontext.Context) {
	panic("test_sentry_error")
}

func innerLibWithoutGoErrors() error {
	return default_errors.New("inner lib error with builtin errors package")
}

func innerLibWithGoErrors() error {
	return errors.SentinelWithStack(ErrInvalid)
}

func WithoutGoErrors(c *webcontext.Context) {
	err := innerLibWithoutGoErrors()
	if err != nil {
		panic(default_errors.Wrap(err, "error calling inner lib"))
	}
	c.JSON(200, "ok")
}

func WithGoErrors(c *webcontext.Context) {
	err := innerLibWithGoErrors()
	if err != nil {
		panic(errors.Wrap(err, "error calling inner lib"))
	}
	c.JSON(200, "ok")
}
