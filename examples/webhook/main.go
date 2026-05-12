// Command webhook is a runnable example of receiving and authenticating
// personal-API webhooks from monobank.
//
// Subscribe a public URL to mono via POST /personal/webhook (see
// PersonalClient.SetWebHook). Every statement event will arrive here as a
// POST /webhook with body and X-Sign/X-Key-Id headers. We fetch the bank's
// signing key once at start-up and refresh it whenever mono rotates it (the
// inbound X-Key-Id stops matching).
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/vtopc/go-monobank"
)

type keyCache struct {
	mu     sync.RWMutex
	client monobank.Client
	key    *monobank.ServerKey
}

func (c *keyCache) refresh(ctx context.Context) error {
	sk, err := c.client.ServerKey(ctx)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.key = sk
	c.mu.Unlock()
	return nil
}

func (c *keyCache) get() *monobank.ServerKey {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.key
}

func main() {
	keys := &keyCache{client: monobank.NewClient(nil)}
	if err := keys.refresh(context.Background()); err != nil {
		log.Fatalf("initial ServerKey: %v", err)
	}
	log.Printf("loaded mono key id=%s", keys.get().ID)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Mono pings the URL with GET when you subscribe a webhook;
			// just answer 200 so the subscription is accepted.
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}

		sk := keys.get()
		if r.Header.Get("X-Key-Id") != sk.ID {
			// Mono rotated its key; refresh and try again.
			if err := keys.refresh(r.Context()); err != nil {
				log.Printf("refresh ServerKey: %v", err)
			}
			sk = keys.get()
		}

		if err := sk.Verify(body, r.Header.Get("X-Sign")); err != nil {
			log.Printf("bad signature from %s: %v", r.RemoteAddr, err)
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}

		event, err := monobank.ParseWebHook(body)
		if err != nil {
			log.Printf("parse webhook: %v", err)
			// Still 200 — payload is authenticated, just unknown.
			return
		}
		t := event.Data.Transaction
		log.Printf("account=%s amount=%d %q hold=%v",
			event.Data.AccountID, t.Amount, t.Description, t.Hold)
	})

	log.Printf("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
