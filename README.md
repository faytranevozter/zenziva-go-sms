# Golang Zenziva SMS Library
Zenziva SMS Online Gateway Library based on Zenziva [Documentation](https://www.zenziva.id/dokumentasi/) with golang.


## Installation
Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:
```bash
$ go get -u github.com/faytranevozter/zenziva-go-sms
```
Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

## Usage
**Reguler type**
```go
import zen "github.com/faytranevozter/zenziva-go-sms"

sms := zen.Zenziva{
  Username: "userkey",
  Password: "passkey",
}
```
**Masking type**
```go
sms := zen.Zenziva{
  Username: "userkey",
  Password: "passkey",
  Type:     "masking",
}
```
**SMS Center type**
```go
sms := zen.Zenziva{
  Username:  "userkey",
  Password:  "passkey",
  Type:      "sms_center",
  Subdomain: "mysubdomain",
}
```
Avalilable type: reguler, masking, sms_center, whatsapp_reguler (coming soon), whatsapp_center (coming soon)

### Sending SMS
**Chaining method**
```go
res, err := sms.To("089765432123").Message("Helaw!").OTP(true).Send()
if err != nil {
  fmt.Println("Failed:", err)
}
```
**Simple method**
```go
res, err := sms.SimpleSend("089765432123", "Helaw!")
if err != nil {
  fmt.Println("Failed:", err)
}
```
**Simple method + otp**
```go
res, err := sms.SimpleSendOTP("089765432123", "Helaw!", true)
if err != nil {
  fmt.Println("Failed:", err)
}
```

### Handling Response
#### Checking SMS Status
```go
res, err := sms.SimpleSend("089765432123", "Helaw!")
if err != nil {
  fmt.Println("Failed:", err)
}

if res.Status {
  // sent
} else {
  // failed
}

```
#### Get Response or Error
```go
// print(res)
// {
//   "message_id": 41,
//   "to": "081111111111",
//   "status": true,
//   "message": "Success"
// }
```

## Credits and License
### Author
Fahrur Rifai [fahrur.dev](https://www.fahrur.dev)  
Twitter [@faytranevozter](https://twitter.com/faytranevozter)

### License
MIT License
