# wesender-go

Officiële Go SDK voor de [WeSender](https://wesender.nl) e-mail API.

## Installatie

```bash
go get github.com/nljerry/wesender-go
```

## Gebruik

```go
package main

import (
    "fmt"
    wesender "github.com/nljerry/wesender-go"
)

func main() {
    ws := wesender.New(os.Getenv("WS_API_KEY"))

    // E-mail versturen
    result, err := ws.SendEmail(wesender.SendEmailInput{
        From:    "noreply@joudomein.nl",
        To:      []string{"klant@voorbeeld.nl"},
        Subject: "Welkom!",
        HTML:    "<h1>Hallo!</h1>",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(result.ID)

    // Domeinen bekijken
    domains, err := ws.ListDomains()
    fmt.Println(domains)
}
```

## Vereisten

Go 1.21+

## Links

- [Documentatie](https://wesender.nl/docs/sdks/go)
- [Issues](https://github.com/nljerry/wesender-go/issues)
