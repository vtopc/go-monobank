// Command webhook is a minimal recipe for receiving and authenticating
// personal-API webhooks from monobank.
//
// Subscribe a public URL via PersonalClient.SetWebHook and every statement
// event will arrive here as POST /webhook with body and X-Sign / X-Key-Id
// headers. This example loads the bank's signing key once at start-up;
// production code should refresh it whenever an incoming X-Key-Id stops
// matching ServerKey.ID (mono's rotation signal).
package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/vtopc/go-monobank"
)

func main() {
	client := monobank.NewClient(nil)
	sk, err := client.ServerKey(context.Background())
	if err != nil {
		log.Fatalf("ServerKey: %v", err)
	}
	log.Printf("loaded mono key id=%s", sk.ID)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Mono pings the URL with GET when subscribing a webhook.
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}
		if err := sk.Verify(body, r.Header.Get("X-Sign")); err != nil {
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		event, err := monobank.ParseWebHook(body)
		if err != nil {
			log.Printf("parse webhook: %v", err)
			return
		}
		t := event.Data.Transaction
		log.Printf("account=%s amount=%d %q hold=%v",
			event.Data.AccountID, t.Amount, t.Description, t.Hold)
	})

	log.Printf("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
