package smsreg
import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"net"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
//

// Stop signal for stoppable listener
type StopListenError struct {}
func (e StopListenError) Error() string {
	return "Stop listening"
}

var stopErr = StopListenError{}

// The tcp keep alive listener that can be stopped

type stoppableTCPListener struct {
	*net.TCPListener
	done <-chan interface {}
}

func NewStoppableTCPListener(addr string, done <-chan interface {}) (net.Listener, error) {

	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return stoppableTCPListener{ln.(*net.TCPListener), done}, err
}

func (ln stoppableTCPListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}

	select {
	case <-ln.done:
		err = stopErr
		return
	default:
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

///////////////////////////////////////////////////////////////////////////////

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
		Addr			: fmt.Sprintf(":%d", port),
		ReadTimeout		: 10 * time.Second,
		WriteTimeout	: 10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	log.Notice("Listening port: %d for incomming sms requests", port)

	go server.ListenAndServe()
}
