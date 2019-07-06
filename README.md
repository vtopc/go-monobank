# go-monobank
Monobank REST API client

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
https://api.monobank.ua/docs/
