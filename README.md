# go-monobank

[![Godoc Reference][godoc-img]][godoc]

Monobank REST API client.
Currently supported only personal authorization(with Token).

## Usage
```go
import (
    "github.com/vtopc/go-monobank"
)

func main(){
    client := monobank.New(nil, monobank.NewPersonalAuthorizer(os.Getenv("TOKEN")))
    response, err := client.ClientInfo()
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
- Linter
- CI for tests

[godoc]: https://godoc.org/github.com/vtopc/go-monobank
[godoc-img]: https://godoc.org/github.com/vtopc/go-monobank?status.svg
