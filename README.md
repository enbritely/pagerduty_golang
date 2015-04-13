# Golang Pagerduty library

```go
package main

import (
    "fmt"
    pagerduty "github.com/enbritely/pagerduty_golang"
)

func main() {
	pd := pagerduty.Pagerduty{
		Servicekey: "my-service-key",
		Client:     "test",
		ClientUrl:  "monitoring.com",
	}
	key, err := pd.CreateIncident("client_api/ourclient", "API endpoint of our client is not operational!", 5)
	if err == nil {
		fmt.Println(key)
	} else {
		fmt.Println(err)
	}

}
```
