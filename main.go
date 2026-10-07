// Command identity is a demo service: it answers /verify and sits at the
// bottom of the call chain (orders -> payments -> identity).
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mcafeelabs/platform/pkg/servicekit"
)

// Config is config-values.schema.yaml.
type Config struct {
	Issuer     string `json:"issuer"`
	SigningKey string `json:"signingKey"`
}

func handler(k *servicekit.Kit, cfg Config) {
	k.Mux.HandleFunc("GET /verify", func(w http.ResponseWriter, r *http.Request) {
		info := k.Info(r.Context())
		info["issuer"] = cfg.Issuer
		// Prove the secret arrived without echoing it.
		info["signingKeyLoaded"] = cfg.SigningKey != ""
		k.JSON(w, http.StatusOK, info)
	})
}

func main() {
	k, err := servicekit.New()
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	if err := k.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	handler(k, cfg)
	if err := k.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
