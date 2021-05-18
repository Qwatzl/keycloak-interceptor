package token

import (
	"io/ioutil"
	"keycloak-interceptor/oidc/shared"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_newRequest(t *testing.T) {
	type args struct {
		tokenEndpoint string
		clientId      string
		code          string
		redirectURI   string
		credentials   *shared.Credentials
	}
	tests := []struct {
		name string
		args args
		want *request
	}{
		{"should create a request without credentials", args{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", nil}, &request{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", nil}},
		{"should create a request with credentials", args{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", &shared.Credentials{Secret: "secret"}}, &request{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", &shared.Credentials{Secret: "secret"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRequest(tt.args.tokenEndpoint, tt.args.clientId, tt.args.code, tt.args.redirectURI, tt.args.credentials); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_toRequest(t1 *testing.T) {
	type fields struct {
		tokenEndpoint string
		clientId      string
		code          string
		redirectURI   string
		credentials   *shared.Credentials
	}

	public := http.Header{}
	confidential := http.Header{}
	public.Set("content-type", "application/x-www-form-urlencoded")
	confidential.Set("content-type", "application/x-www-form-urlencoded")
	confidential.Set("authorization", "basic dGVzdC1jbGllbnQ6c2VjcmV0")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", "code123")
	form.Set("redirect_uri", "http://localhost:4000/redirect")

	tests := []struct {
		name    string
		fields  fields
		want    want
		wantErr bool
	}{
		{"should create a request without credentials", fields{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", nil}, want{http.MethodPost, "localhost:8080", "/token", public, form.Encode()}, false},
		{"should create a request with credentials", fields{"http://localhost:8080/token", "test-client", "code123", "http://localhost:4000/redirect", &shared.Credentials{Secret: "secret"}}, want{http.MethodPost, "localhost:8080", "/token", confidential, form.Encode()}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &request{
				tokenEndpoint: tt.fields.tokenEndpoint,
				clientId:      tt.fields.clientId,

				code:        tt.fields.code,
				redirectURI: tt.fields.redirectURI,
				credentials: tt.fields.credentials,
			}
			got, err := t.toRequest()
			if (err != nil) != tt.wantErr {
				t1.Errorf("toRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareHttpRequests(got, tt.want, t1) {

			}
		})
	}
}

type want struct {
	method string
	host   string
	path   string
	header http.Header
	body   string
}

func compareHttpRequests(got *http.Request, want want, t *testing.T) bool {
	if got.Method != want.method {
		t.Errorf("toRequest() got = '%v', want '%v'", got.Method, want.method)
	}
	if got.URL.Host != want.host {
		t.Errorf("toRequest() got = '%v', want '%v'", got.URL.Host, want.host)
	}
	if got.URL.Path != want.path {
		t.Errorf("toRequest() got = '%v', want '%v'", got.URL.Path, want.path)
	}
	if !reflect.DeepEqual(got.Header, want.header) {
		t.Errorf("toRequest() got = %v, want %v", got.Header, want.header)
	}

	gotBody, _ := ioutil.ReadAll(got.Body)
	if string(gotBody) != want.body {
		t.Errorf("toRequest() got = %v, want %v", gotBody, want.body)
	}

	return true
}
