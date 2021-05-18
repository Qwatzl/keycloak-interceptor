package oidc

import (
	"keycloak-interceptor/oidc/authentication"
	"keycloak-interceptor/oidc/http2"
	"keycloak-interceptor/oidc/shared"
	"keycloak-interceptor/oidc/token"
	"net/http"
)

const (
	oidcCode = "code"
	empty    = ""
)

type Keycloak struct {
	keycloakUIDCJsonPath string
	config               *shared.JSON
	endpoints            *shared.Endpoints
	tokenVerifier        *token.Verifier
}

func Init(keycloakJSON string) *Keycloak {
	k := &Keycloak{
		keycloakUIDCJsonPath: keycloakJSON,
		config:               shared.LoadConfiguration(keycloakJSON),
	}

	k.endpoints = shared.InitEndpoints(k.config)
	k.tokenVerifier = token.InitVerifier(k.endpoints, k.config)

	return k
}

// The Authorization Code Flow goes through the following steps.
//   - Client prepares an Authentication request containing the desired request parameters.
//   - Client sends the request to the Authorization Server.
//   - Authorization Server Authenticates the End-User.
//   - Authorization Server obtains End-User Consent/Authorization.
//   - Authorization Server sends the End-User back to the Client with an Authorization Code.
//   - Client requests a response using the Authorization Code at the Token Endpoint.
//   - Client receives a response that contains an ID Token and Access Token in the response body.
//   - Client validates the ID verifier and retrieves the End-User's Subject Identifier.
func (k *Keycloak) Intercept(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// receive and validate code
		queryParams := r.URL.Query()
		code := queryParams.Get(oidcCode)
		if code != empty && k.tokenVerifier.Verify(code, toRedirectURI(r)) {
			next.ServeHTTP(w, r)
			return
		}

		// redirect to login page
		authReq := authentication.NewRequest(k.endpoints.Authorization, k.config.Resource, toRedirectURI(r))
		w.Header().Set(http2.HeaderLocation, authReq.ToString())
		w.WriteHeader(http.StatusFound)
	})
}

func toRedirectURI(r *http.Request) string {
	return "http://" + r.Host + r.URL.Path
}
