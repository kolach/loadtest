package smsreg

import (
	"time"
	"net/http"
)

type AsyncSmsReq struct {
	url 		string				// URI of remote service
	respChan 	chan<- *SmsIn 		// sms response will be send to this channel
	timestamp 	time.Time			// time the request was issued
	*SmsOut							// request body
}

func NewAsyncSmsReq(url string, smsOut *SmsOut) (*AsyncSmsReq, <-chan *SmsIn) {
	respChan := make(chan *SmsIn)
	return &AsyncSmsReq{
		url			: url,
		SmsOut		: smsOut,
		respChan	: respChan,
		timestamp	: time.Now(),
	}, respChan
}

func (r *AsyncSmsReq) Cancel() {
	r.RespondWith(nil)
}

func (r *AsyncSmsReq) RespondWith(smsIn *SmsIn) {
	r.respChan <- smsIn
}

func (r *AsyncSmsReq) IsLonger(d time.Duration) bool {
	return time.Since(r.timestamp) > d
}

func (r *AsyncSmsReq) Timestamp() {
	r.timestamp = time.Now()
}

func (r *AsyncSmsReq) HttpRequest() (*http.Request, error)  {
	return NewSmsHttpRequest(r.url, r.SmsOut)
}

