package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcafeelabs/platform/pkg/servicekit"
)

func TestVerify(t *testing.T) {
	t.Setenv("SERVICE_NAME", "identity")
	k, err := servicekit.New()
	if err != nil {
		t.Fatal(err)
	}
	handler(k, Config{Issuer: "https://id", SigningKey: "k"})
	srv := httptest.NewServer(k.Handler())
	defer srv.Close()
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/verify", nil)
	req.Header.Set("baggage", "sandbox=demo")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&got)
	if got["service"] != "identity" || got["requestSandbox"] != "demo" || got["signingKeyLoaded"] != true {
		t.Fatalf("got %v", got)
	}
}
