# go-monobank

[![Godoc Reference][godoc-img]][godoc]

Monobank REST API client.
Currently, supported only personal authorization(with Token).

## Usage
```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/vtopc/go-monobank"
)

func main(){
    client := monobank.New(nil).WithAuth(monobank.NewPersonalAuthorizer(os.Getenv("TOKEN")))
    response, _ := client.ClientInfo(context.Background())
    fmt.Println(response)
}
```

## Official docs
 - https://api.monobank.ua/docs/

## Similar projects
- https://github.com/shal/mono
- https://github.com/artemrys/go-monobank-api

## TODO
- Corporate Authorization
- Corporate API(init/check auth)
- More unit tests
- CI for tests

[godoc]: https://godoc.org/github.com/vtopc/go-monobank
[godoc-img]: https://godoc.org/github.com/vtopc/go-monobank?status.svg
