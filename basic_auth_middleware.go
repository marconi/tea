package tea

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var basicAuthError = errors.New("Invalid authentication")

type basicAuthUserIdCtxKey struct{}
type authenticator func(string, string) bool
type authorizator func(string, *http.Request) bool

type basicAuthMiddleware struct {
	realm         string
	authFunc      authenticator
	authorizeFunc authorizator
	next          http.Handler
}

func (m *basicAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		m.unauthorized(w)
		return
	}

	userId, password, err := m.decodeAuthHeader(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(m.authFunc(userId, password) && m.authorizeFunc(userId, r)) {
		m.unauthorized(w)
		return
	}

	c := context.WithValue(r.Context(), basicAuthUserIdCtxKey{}, userId)
	m.next.ServeHTTP(w, r.WithContext(c))
}

func (m *basicAuthMiddleware) decodeAuthHeader(header string) (string, string, error) {
	parts := strings.Split(header, " ")
	if !(len(parts) == 2 && parts[0] == "Basic") {
		return "", "", basicAuthError
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", basicAuthError
	}

	creds := strings.Split(string(decoded), ":")
	if len(creds) != 2 {
		return "", "", basicAuthError
	}

	return creds[0], creds[1], nil
}

func (m *basicAuthMiddleware) unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%s", m.realm))
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func NewBasicAuthMiddleware(realm string, authFunc authenticator, authorizeFunc authorizator) Middleware {
	return func(h http.Handler) http.Handler {
		bam := &basicAuthMiddleware{realm, authFunc, authorizeFunc, h}
		if bam.authorizeFunc == nil {
			bam.authorizeFunc = func(userId string, r *http.Request) bool {
				return true
			}
		}
		return bam
	}
}
