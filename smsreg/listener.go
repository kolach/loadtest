package smsreg
import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type apiStatResp struct {
	ResultCode int `json:"resultCode"`
	Message string `json:"message"`
}

// Listen incoming SMS
// port - port to listen
// done - channel to listen for done event - on receive should stop listener server
// smsIns - array of channels (acts like topic) to redirect incoming messages
func Listen(port int, done <-chan struct{}, smsIns... chan<- *SmsIn) {

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
			for _, smsIn := range smsIns {
				smsIn <- msg
			}
		}()

		// and responding with ApiStatusResponse structure that server expects to get
		w.Header().Set("Content-Type", "application/json")
		json, _ := json.Marshal(&apiStatResp{200, "OK"})
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
}
