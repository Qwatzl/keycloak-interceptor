package main

import (
	"keycloak-interceptor/oidc"
	"log"
	"net/http"
)

func serviceMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(serviceMock)
	kInterceptor := oidc.Init("keycloak.json")
	mux.Handle("/users", kInterceptor.Intercept(finalHandler))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
