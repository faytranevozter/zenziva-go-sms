package zenziva

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

// Zenziva model
type Zenziva struct {
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Type      string  `json:"type"`
	Subdomain string  `json:"subdomain"`
	Payload   payload `json:"payload"`
}

type payload struct {
	Username string `json:"username" url:"userkey"`
	Password string `json:"password" url:"passkey"`
	To       string `json:"to" url:"nohp"`
	Message  string `json:"message" url:"pesan"`
	OTP      bool   `json:"otp,omitempty" url:"-"`
}

type zenResponse map[string]interface{}

// Response API
type Response struct {
	MessageID int64  `json:"message_id"`
	To        string `json:"to"`
	Status    bool   `json:"status"`
	Message   string `json:"message"`
}

// static const
const (
	SCHEME = "https"
	DOMAIN = "zenziva.net"
)

var availableType []string = []string{"reguler", "masking", "sms_center", "whatsapp_reguler", "whatsapp_center"}

// To : set phone number
func (zen *Zenziva) To(phone string) *Zenziva {
	zen.Payload.To = phone
	return zen
}

// Message : set message
func (zen *Zenziva) Message(text string) *Zenziva {
	zen.Payload.Message = text
	return zen
}

// OTP : set message is otp or not
func (zen *Zenziva) OTP(status bool) *Zenziva {
	zen.Payload.OTP = status
	return zen
}

// Send : Used to send message
func (zen *Zenziva) Send() (res Response, err error) {
	err = zen.initRequest()
	if err != nil {
		return
	}

	zen.assignConfig()

	res, err = zen.requestAPI(zen.Payload)
	if err != nil {
		return
	}

	return
}

// SimpleSend : Used to send message with simpler method
func (zen *Zenziva) SimpleSend(to, message string) (res Response, err error) {
	return zen.To(to).Message(message).Send()
}

// SimpleSendOTP : Used to send message with simpler method + otp
func (zen *Zenziva) SimpleSendOTP(to, message string, otp bool) (res Response, err error) {
	return zen.To(to).Message(message).OTP(otp).Send()
}

func (zen *Zenziva) typeSMS(smsType string) {
	zen.Type = smsType
}

func (zen *Zenziva) subdomain(subdomain string) {
	zen.Subdomain = subdomain
}

func (zen *Zenziva) initRequest() error {
	if zen.Username == "" {
		return errors.New("Username is requried")
	}

	if zen.Password == "" {
		return errors.New("Password is requried")
	}

	if zen.Type == "" {
		zen.Type = "reguler"
	}

	if !inArray(zen.Type, availableType) {
		return errors.New("Type/Package should one of (" + strings.Join(availableType, ", ") + ")")
	}

	if inArray(zen.Type, []string{"sms_center", "whatsapp_center"}) && zen.Subdomain == "" {
		return errors.New("Subdomain is requried if type sms_center or whatsapp_center")
	}

	return nil
}

func (zen *Zenziva) assignConfig() *Zenziva {
	zen.Payload.Username = zen.Username
	zen.Payload.Password = zen.Password
	return zen
}

func (zen *Zenziva) getSendSMSURL() string {
	if inArray(zen.Type, []string{"reguler", "whatsapp_reguler"}) {
		zen.Subdomain = "gsm"
	} else if zen.Type == "masking" {
		zen.Subdomain = "masking"
	}

	path := ""
	if inArray(zen.Type, []string{"reguler", "masking", "sms_center"}) {
		if zen.Payload.OTP {
			path = "api/sendOTP/"
		} else {
			path = "api/sendsms/"
		}
	} else if zen.Type == "whatsapp_reguler" {
		path = "api/sendWA/"
	} else if zen.Type == "whatsapp_center" {
		path = "api/WAsendMsg/"
	}

	url := SCHEME + "://" + zen.Subdomain + "." + DOMAIN + "/" + path

	return url
}

func (zen *Zenziva) requestAPI(data payload) (resAPI Response, err error) {
	url := zen.getSendSMSURL()

	v, _ := query.Values(data)
	payload := strings.NewReader(v.Encode())
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	zenRes := zenResponse{}
	err = json.Unmarshal(body, &zenRes)
	if err != nil {
		return
	}

	messageID := 0
	if reflect.TypeOf(zenRes["messageId"]).String() == "string" {
		if zenRes["messageId"] != "" {
			messageID, _ = strconv.Atoi(zenRes["messageId"].(string))
		}
	} else if reflect.TypeOf(zenRes["messageId"]).String() == "float64" {
		messageID = int(zenRes["messageId"].(float64))
	}

	resAPI = Response{
		MessageID: int64(messageID),
		To:        zenRes["to"].(string),
		Status:    zenRes["status"].(string) == "1",
		Message:   zenRes["text"].(string),
	}

	return resAPI, nil
}
