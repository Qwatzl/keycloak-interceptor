package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keycloak-interceptor/oidc/http2"
	"log"
	"os"
)

type Credentials struct {
	Secret string `json:"secret"`
}

type JSON struct {
	Realm                   string      `json:"realm"`
	AuthServerURL           string      `json:"auth-server-url"`
	SSLRequired             string      `json:"ssl-required"`
	Resource                string      `json:"resource"`
	VerifyTokenAudience     bool        `json:"verify-verifier-audience"`
	Credentials             Credentials `json:"credentials"`
	UseResourceRoleMappings bool        `json:"use-resource-role-mappings"`
	ConfidentialPort        int         `json:"confidential-port"`
}

type Endpoints struct {
	Authorization string `json:"authorization_endpoint"`
	Token         string `json:"token_endpoint"`
	UserInfo      string `json:"userinfo_endpoint"`
	EndSession    string `json:"end_session_endpoint"`
}

func LoadConfiguration(uidcJSONPath string) *JSON {
	jsonFile, err := os.Open(uidcJSONPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var config JSON
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func InitEndpoints(config *JSON) *Endpoints {
	if config == nil {
		log.Fatal("auth server initialization failed")
	}

	url := fmt.Sprintf("%srealms/%s/.well-known/openid-configuration", config.AuthServerURL, config.Realm)
	var endpoints Endpoints
	if err := http2.Request(url, &endpoints); err != nil {
		log.Fatal("auth server initialization failed")
	}

	return &endpoints
}
