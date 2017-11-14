# go-dexcom
A go client for the [Dexcom](https://developer.dexcom.com/getting-started) API


## Usage

```sh
go get github.com/healthimation/go-dexcom/dexcom
```

```golang
import (
    "log"
    "time"

    "github.com/healthimation/go-dexcom/dexcom"
)

func main() {
    timeout := 5 * time.Second
    client := dexcom.NewClient("my dexcom client id", "my dexcom client secret", timeout)

    
    userToken, err := client.GetUser()(context.Background(), userAuthCode, redirectURI) 
    if err != nil {
        log.Printf("Error fetching user token: %s", err.Error())
    }

    //Get EGVs
    egvs, err := client.GetEGVs(context.Background(), userToken.AccessToken)
    if err != nil {
        log.Printf("Error fetching user egvs: %s", err.Error())
    }
}
```
