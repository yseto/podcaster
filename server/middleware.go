package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/yseto/podcaster/ent"
)

func CreateMiddlewareEmptyContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(emptyUserContext(r.Context())))
		})
	}
}

func CreateMiddleware(ent *ent.Client) (func(next http.Handler) http.Handler, error) {
	spec, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: newAuthenticator(ent),
			},
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				if statusCode == 401 && strings.HasPrefix(message, "security requirements failed") {
					w.Header().Add("WWW-Authenticate", "Basic realm=podcaster")
				}
				w.WriteHeader(statusCode)
			},
		})

	return validator, nil
}
