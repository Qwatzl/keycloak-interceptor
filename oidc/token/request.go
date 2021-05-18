package token

import (
	"encoding/base64"
	"keycloak-interceptor/oidc/http2"
	"keycloak-interceptor/oidc/shared"
	"net/http"
	"net/url"
	"strings"
)

// From: https://openid.net/specs/openid-connect-core-1_0.html#TokenRequest
// A Client makes a Token request by presenting its Authorization Grant (in the form of an Authorization Code) to the Token Endpoint using the grant_type value authorization_code, as described in Section 4.1.3 of OAuth 2.0 [RFC6749]. If the Client is a Confidential Client, then it MUST authenticate to the Token Endpoint using the authentication method registered for its client_id, as described in Section 9.
// The Client sends the parameters to the Token Endpoint using the HTTP POST method and the Form Serialization, per Section 13.2, as described in Section 4.1.3 of OAuth 2.0 [RFC6749].
// The following is a non-normative example of a Token request (with line wraps within values for display purposes only):
//
//  POST /token HTTP/1.1
//  Host: server.example.com
//  Content-Type: application/x-www-form-urlencoded
//  Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
//
//  grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA
//    &redirect_uri=https%3A%2F%2Fclient.example.org%2Fcb
type request struct {
	tokenEndpoint string
	clientId      string
	code          string
	redirectURI   string
	credentials   *shared.Credentials
}

func newRequest(tokenEndpoint string, clientId string, code string, redirectURI string, credentials *shared.Credentials) *request {
	return &request{
		tokenEndpoint: tokenEndpoint,
		clientId:      clientId,
		code:          code,
		redirectURI:   redirectURI,
		credentials:   credentials,
	}
}

func (t *request) toRequest() (*http.Request, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", t.code)
	form.Set("redirect_uri", t.redirectURI)

	req, err := http.NewRequest(http.MethodPost, t.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set(http2.HeaderContentType, http2.HeaderContentTypeFormURLEncoded)
	if t.credentials != nil {
		encoder := base64.StdEncoding
		cred := t.clientId + ":" + t.credentials.Secret
		req.Header.Set(http2.HeaderAuthorization, "basic "+encoder.EncodeToString([]byte(cred)))
	}

	return req, nil
}
