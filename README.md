# go-monobank

[![Godoc Reference][godoc-img]][godoc-url] [![CI][ci-img]][ci-url] [![codecov][codecov-img]][codecov-url]

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

[godoc-img]: https://godoc.org/github.com/vtopc/go-monobank?status.svg
[godoc-url]: https://godoc.org/github.com/vtopc/go-monobank

[ci-img]: https://github.com/vtopc/go-monobank/workflows/CI/badge.svg
[ci-url]: https://github.com/vtopc/go-monobank/actions?query=workflow%3A%22CI%22

[codecov-img]: https://codecov.io/gh/vtopc/go-monobank/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/vtopc/go-monobank
