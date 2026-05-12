# go-monobank

[![Godoc Reference][godoc-img]][godoc-url] [![CI][ci-img]][ci-url] [![codecov][codecov-img]][codecov-url]

Monobank REST API client.

## Features
- Personal API(with Token authorization).
- API for providers(corporate) with authorization.
- Webhooks(including API for providers).
- Jars(only in Personal API).
- Webhook signature verification (ECDSA secp256k1) via the bank's public key.

## Installation
```shell
go get github.com/vtopc/go-monobank@latest
```
This will update yours go.mod file.

## Usage examples
NOTE: Do not forget to check errors.

#### Public client
```go
package main

import (
    "context"
    "fmt"

    "github.com/vtopc/go-monobank"
)

func main() {
    // Create public client.
    client := monobank.NewClient(nil)

    response, _ := client.Currency(context.Background())
    fmt.Println(response)
}
```

#### Personal client
```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/vtopc/go-monobank"
)

func main() {
    token := os.Getenv("TOKEN")

    // Create authorized client.
    client := monobank.NewPersonalClient(nil).WithAuth(monobank.NewPersonalAuthorizer(token))

    response, _ := client.ClientInfo(context.Background())
    fmt.Println(response)
}
```

#### Corporate client
```go
package main

import (
    "context"
    "fmt"

    "github.com/vtopc/go-monobank"
)

var secKey []byte // put here you private key

const webhook = "https://example.com/webhook"

func main() {
    // Create auth creator.
    authMaker, _ := monobank.NewCorpAuthMaker(secKey)

    // Create authorized client.
    client, _ := monobank.NewCorporateClient(nil, authMaker)

    // If the user is not authorized yet, do next:
    resp, _ := client.Auth(context.Background(), webhook, monobank.PermSt, monobank.PermPI)

    // Send `resp.AcceptURL` to the user and wait until it will authorize your client
    // in Monobank app on mobile, you will get GET request on `webhook` when it will be done.
    // See Documentation for details.
    // Store `resp.RequestID` somewhere.
    requestID := resp.RequestID

    // If user authorized already:
    response, _ := client.ClientInfo(context.Background(), requestID)
    fmt.Println(response)
}
```

#### Verifying webhook signatures
Mono signs each webhook with ECDSA over secp256k1. Fetch the bank's public key
from `/bank/sync` once and use it to verify every incoming request:
```go
package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "sync"

    "github.com/vtopc/go-monobank"
)

func main() {
    client := monobank.NewClient(nil)

    var (
        mu  sync.RWMutex
        key *monobank.ServerKey
    )
    refresh := func() error {
        sk, err := client.ServerKey(context.Background())
        if err != nil {
            return err
        }
        mu.Lock()
        key = sk
        mu.Unlock()
        return nil
    }
    if err := refresh(); err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        body, _ := io.ReadAll(r.Body)

        mu.RLock()
        k := key
        mu.RUnlock()
        // If X-Key-Id no longer matches, mono rotated the key — refresh and retry.
        if r.Header.Get("X-Key-Id") != k.ID {
            if err := refresh(); err == nil {
                mu.RLock()
                k = key
                mu.RUnlock()
            }
        }

        if err := monobank.VerifyWebhookSignature(k.PubKey, body, r.Header.Get("X-Sign")); err != nil {
            http.Error(w, "bad signature", http.StatusUnauthorized)
            return
        }
        // ... process body as monobank.WebHookResponse
        w.WriteHeader(http.StatusOK)
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Documentation
- [Official](https://api.monobank.ua/docs/)
- [Unofficial](https://gist.github.com/Sominemo/8714a82e26a268c30e4a332b0b2fd943)

## Similar projects
- https://github.com/shal/mono (last update 16.05.2020)
- https://github.com/artemrys/go-monobank-api (no corporate API)

## TODO
- More unit tests

[godoc-img]: https://godoc.org/github.com/vtopc/go-monobank?status.svg
[godoc-url]: https://godoc.org/github.com/vtopc/go-monobank

[ci-img]: https://github.com/vtopc/go-monobank/workflows/CI/badge.svg
[ci-url]: https://github.com/vtopc/go-monobank/actions?query=workflow%3A%22CI%22

[codecov-img]: https://codecov.io/gh/vtopc/go-monobank/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/vtopc/go-monobank
