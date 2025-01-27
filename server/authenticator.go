package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/yseto/podcaster/ent"
	"github.com/yseto/podcaster/ent/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	errNoAuthHeader           = errors.New("authorization header is missing")
	errInvalidAuthHeader      = errors.New("authorization header is malformed")
	errNotImplemented         = errors.New("not implemted")
	errGettingHeaderParameter = errors.New("getting header param")
	errUserNotFound           = errors.New("user not found")
)

func newAuthenticator(ent *ent.Client) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return authenticate(ctx, ent, input)
	}
}

func authorizationFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", errNoAuthHeader
	}
	prefix := "Basic "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", errInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func authenticate(ctx context.Context, ent *ent.Client, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "basicAuth" {
		return errNotImplemented
	}

	basicBase64enc, err := authorizationFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		fmt.Println(err)
		return errGettingHeaderParameter
	}

	// fmt.Println(basicBase64enc)

	var basicAuth []byte
	basicAuth, err = base64.StdEncoding.DecodeString(basicBase64enc)
	if err != nil {
		fmt.Println(err)
		return errGettingHeaderParameter
	}

	param := strings.SplitN(string(basicAuth), ":", 2)

	// fmt.Println(param)

	u, err := ent.Users.Query().Where(users.NameEQ(param[0])).First(ctx)
	if err != nil {
		return errUserNotFound
	}

	if check := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(param[1])); check != nil {
		return errUserNotFound
	}

	newUserContext(ctx, uint64(u.ID))

	return nil
}
