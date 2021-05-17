package token

import (
	"io/ioutil"
	"keycloak-interceptor/oidc/shared"
	"log"
	"net/http"
)

type Verifier struct {
	client *http.Client
	e      *shared.Endpoints
	c      *shared.JSON
}

func InitVerifier(endpoints *shared.Endpoints, config *shared.JSON) *Verifier {
	return &Verifier{
		client: http.DefaultClient,
		e:      endpoints,
		c:      config,
	}
}

func (v *Verifier) Verify(code string, redirectURI string) bool {
	tr := newRequest(v.e.Token, v.c.Resource, code, redirectURI, &v.c.Credentials)
	req, err := tr.toRequest()

	if err != nil {
		log.Println("unable to create a request oidc token request", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("unable to request an oidc token", err)
		return false
	}

	log.Printf("oidc token request responded with status code: %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("oidc token request response: %v\n", string(body))
	}
	return resp.StatusCode == http.StatusOK
}
