package smsreg

import (
	"bytes"
	"net/http"
	"encoding/json"
)

type SmsOut struct {
	Originator 	string `json:"originator"`
	Receiver 	string `json:"receiver"`
	Text 		string `json:"text"`
}

func NewSmsOut(phone, text string) *SmsOut {
	return &SmsOut{
		Originator: phone,
		Text: text,
		Receiver: "9999999",
	}
}

func NewSmsHttpRequest(url string, sms *SmsOut) (*http.Request, error) {
	json, err := json.Marshal(sms)

	if err != nil { return nil, err }

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	if err != nil { return nil, err }

	httpReq.Header.Set("Content-Type", APP_JSON)
	httpReq.Header.Set("Accept", APP_JSON)
	return httpReq, nil
}
