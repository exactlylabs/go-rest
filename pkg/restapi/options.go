package restapi

import (
	"github.com/gorilla/handlers"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var defaultPropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

type ServerOption func(server *WebServer) error

func WithRecoveryMiddleware(middleware Middleware) ServerOption {
	return func(server *WebServer) error {
		server.recoveryMiddleware = middleware
		return nil
	}
}

func WithRequestLoggerMiddleware(middleware Middleware) ServerOption {
	return func(server *WebServer) error {
		server.loggerMiddleware = middleware
		return nil
	}
}

// WithTracing option adds an OpenTelemetry's tracing SDK implementation to the server
func WithTracing(traceName string, tp trace.TracerProvider) ServerOption {
	return func(server *WebServer) error {
		server.tracerProvider = tp
		if server.propagator == nil {
			server.propagator = defaultPropagator
		}
		server.baseCtx.Tracer = server.tracerProvider.Tracer(traceName)
		return nil
	}
}

// WithTracePropagator replaces the default propagator with the one provided
func WithTracePropagator(propagator propagation.TextMapPropagator) ServerOption {
	return func(server *WebServer) error {
		server.propagator = propagator
		return nil
	}
}

type CORSOptions struct {
	Origins []string
	Methods []string
	Headers []string
}

var defaultCORSOptions = CORSOptions{
	Origins: nil,
	Methods: []string{"GET", "POST", "PUT", ""},
}

func WithCORS(opts CORSOptions) ServerOption {
	corsOpts := []handlers.CORSOption{}
	if opts.Origins != nil {
		corsOpts = append(corsOpts, handlers.AllowedOrigins(opts.Origins))
	}
	if opts.Headers != nil {
		corsOpts = append(corsOpts, handlers.AllowedHeaders(opts.Headers))
	}
	if opts.Methods != nil {
		corsOpts = append(corsOpts, handlers.AllowedMethods(opts.Methods))
	}
	return func(server *WebServer) error {
		server.AddMiddlewares(handlers.CORS(corsOpts...))
		return nil
	}
}
