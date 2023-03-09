package auth

import (
	"api_gateway/pkg/middlewares/base"
	"api_gateway/pkg/middlewares/logs"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	goauth "github.com/abbot/go-http-auth"
)

const (
	TypeName     = "BasicAuth"
	defaultRealm = "api-gateway"
)

type basicAuth struct {
	next         http.Handler
	auth         *goauth.BasicAuth
	users        map[string]string
	headerField  string
	removeHeader bool
	name         string
}

type BasicAuth struct {
	base.Config
	Users []string `json:"users,omitempty"`
}

// NewBasic creates a basicAuth middleware.
func NewBasic(ctx context.Context, next http.Handler, authConfig *BasicAuth, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Msg("Creating middleware")

	users, err := getUsers(authConfig.Users, basicUserParser)
	if err != nil {
		return nil, err
	}

	ba := &basicAuth{
		next:         next,
		users:        users,
		removeHeader: true,
		name:         name,
	}

	realm := defaultRealm

	ba.auth = &goauth.BasicAuth{Realm: realm, Secrets: ba.secretBasic}

	return ba, nil
}

func (b *basicAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := logs.GetLogger(req.Context(), b.name, TypeName)

	user, password, ok := req.BasicAuth()
	if ok {
		secret := b.auth.Secrets(user, b.auth.Realm)
		if secret == "" || !goauth.CheckSecret(password, secret) {
			ok = false
		}
	}

	if !ok {
		logger.Debug().Msg("Authentication failed")

		b.auth.RequireAuth(rw, req)
		return
	}

	logger.Debug().Msg("Authentication succeeded")
	req.URL.User = url.User(user)

	if b.headerField != "" {
		req.Header[b.headerField] = []string{user}
	}

	if b.removeHeader {
		logger.Debug().Msg("Removing authorization header")
		req.Header.Del(authorizationHeader)
	}
	b.next.ServeHTTP(rw, req)
}

func (b *basicAuth) secretBasic(user, realm string) string {
	if secret, ok := b.users[user]; ok {
		return secret
	}

	return ""
}

func basicUserParser(user string) (string, string, error) {
	split := strings.Split(user, ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("error parsing BasicUser: %v", user)
	}
	return split[0], split[1], nil
}
