# yeastartgsms
TG SMS API in Go

Usage example:

```go
func main() {
	port := 1
	dest := "02342324232"
	message := "Hello, this is a test message!"
	id := "xxxb3"

	err := SendSMS(port, dest, message, id)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		fmt.Println("SMS sent successfully!")
	}
}
```