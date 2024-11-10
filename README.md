# yeastartgsms
TG SMS API in Go

Usage example:

```go
import (
	"fmt"
	"log"
	yeastartgsms "github.com/andrejro2/yeastartgsms"
)
func main() {
    gatewayHostname := "192.168.1.1:5000"
    gatewayUser := "test"
    gatewaySecret := "test"
	port := 1
	dest := "02342324232"
	message := "Hello, this is a test message!"
	id := "xxxb3"
	// Set connection Data for SMS gateway.
	yeastartgsms.SetConnection(gatewayHostname, gatewayUser, gatewaySecret)
    // Send SMS
	err := yeastartgsms.SendSMS(port, dest, message, id)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		fmt.Println("SMS sent successfully!")
	}
}
```