package smsreg


import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

const (
	URL = "/api/sms/new/transactional/psms"
)


// Incoming Caipirinha Server SMS
type SmsIn struct {
	Sender 			string `json:"sender"`
	Recipient 		string `json:"recipient"`
	Text 			string `json:"text"`
	CreateDateTime 	string `json:"createDateTime"`
	Carrier 		string `json:"carrier"`
}

// convert SmsRes to string
func (sms *SmsIn) String() string {
	// Carrier and Carrier are not interesting fields for current purposes
	return fmt.Sprintf("Recipient: %s, Text: %s, CreateDateTime: %s", sms.Recipient, sms.Text, sms.CreateDateTime)
}


type apiStatusResponse struct {
	ResultCode int `json:"resultCode"`
	Message string `json:"message"`
}

// Listen incoming SMS
// port - port to listen
// out  - incoming output channel
func Listen(port int) <-chan *SmsIn {

	smsIn := make(chan *SmsIn, 10000)

	http.HandleFunc(URL, func (w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		msg := &SmsIn{}
		err := json.Unmarshal(body, msg)
		if err != nil {
			log.Error("Error unmarshaling message: %s", err)
		}

		log.Debug("New message received: %s", msg)

		// async write to output channel
		go func() {
			smsIn <- msg
		}()

		// and responding with ApiStatusResponse structure that server expects to get
		w.Header().Set("Content-Type", "application/json")
		json, _ := json.Marshal(&apiStatusResponse{200, "OK"})
		w.Write(json)

	})

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Additional server params
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	log.Notice("Listening port: %d for incomming sms requests", port)

	go server.ListenAndServe()

	return smsIn
}
