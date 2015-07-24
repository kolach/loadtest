package smsreg

import (
	"github.com/kolach/loadtest/random"
	"net/http"
	"time"
	"runtime"
)

const (
	UNREGISTER_CMD 	= "SAIR"
	REGISTER_CMD 	= "Register me"
	APP_JSON 		= "application/json"
	YES 			= "Si"

	HEALTH_CHECK_INTERVAL = time.Second * 5
	MAX_RESPONSE_WAITTIME = time.Second * 20
	STATS_INTERVAL		  = time.Second * 30
)


type Producer struct {

	url 		string				// URL to send SMS http request
	startTime 	time.Time			// time of production start

	serverRespChan <-chan *SmsIn	// read only channel to listen of server responses

	asyncReqChan chan *AsyncSmsReq	// channel to schedule sms request with async response

	semChan chan int 				// semaphore channel
	count int 						// produce this number of users

	phoneGen 	<-chan string		// random phone number generator
	addressGen 	<-chan string		// random address generator

	countChan chan bool				// channel to count registrations
}


func NewProducer(url string, count, concurrency, port int) *Producer {

	// semaphore channel
	semChan 	 := make(chan int, concurrency)

	// channel to submit sms requests
	asyncReqChan := make(chan *AsyncSmsReq, 1000)

	// channel to count registrations
	countChan := make(chan bool, 2*concurrency)

	// random generators
	phoneGen, err := random.PhoneNumberGen("MX", "55")
	if err != nil { panic(err) }
	// it takes a while to load postal codes database
	addressGen, err := random.AddressGen("MX")
	if err != nil { panic(err) }

	serverRespChan := Listen(port)

	p := &Producer{
		url				: url,
		serverRespChan	: serverRespChan,
		asyncReqChan	: asyncReqChan,
		semChan			: semChan,
		count			: count,
		phoneGen		: phoneGen,
		addressGen		: addressGen,
		countChan		: countChan,
	}

	return p
}


func (p *Producer) logSummary(totalProduced, totalFailures int) {
	log.Notice("TIME ELAPSED: %v, TOTAL PRODUCED: %d, TOTAL FAILURES: %d",
		time.Since(p.startTime), totalProduced, totalFailures)
}

func (p *Producer) logStats(totalProduced, totalFailures, _totalProduced, _totalFailures int) (int, int) {
	log.Notice("Time elapsed: %v, Total Produced: %d (+%d), Total Failures: %d (+%d), Goroutines: %d",
		time.Since(p.startTime), totalProduced, totalProduced - _totalProduced,
		totalFailures, totalFailures - _totalFailures, runtime.NumGoroutine())
	return totalProduced, totalFailures
}

func (p *Producer) handleRequests(done chan<- struct{}, killChan <-chan struct{}) {

	httpClient 		  	:= &http.Client{}
	healthCheckTicker 	:= time.NewTicker(HEALTH_CHECK_INTERVAL)
	statsTicker 		:= time.NewTicker(STATS_INTERVAL)
	requests 		  	:= make(map[string] *AsyncSmsReq)

	totalProduced 	  := 0
	totalFailures 	  := 0
	_totalProduced 	  := 0 // prev
	_totalFailures 	  := 0 // prev

	for {
		select {

		// make request, and register to wait for response
		case asyncSmsReq := <-p.asyncReqChan:

			if oldReq, exists := requests[asyncSmsReq.Originator]; exists {
				log.Critical("Request map already contains a request: %s. Cancelling old one", asyncSmsReq)
				oldReq.Cancel()
			}

			httpRequest, err := asyncSmsReq.HttpRequest()
			if err != nil {
				log.Error("Failed to create http request object, %s", err)
				asyncSmsReq.Cancel()
			} else {
				// TODO parse and analyse http response
				log.Debug("Sending sms http request: %s", asyncSmsReq.SmsOut)
				_, err := httpClient.Do(httpRequest)
				if err != nil {
					log.Error("Failed to make http request, %s", err)
					asyncSmsReq.Cancel()
				} else {
					requests[asyncSmsReq.Originator] = asyncSmsReq
					asyncSmsReq.Timestamp()
				}
			}


		// response received, unregister request from waiting and send response to
		// response channel
		case smsIn := <-p.serverRespChan:
			asyncSmsReq := requests[smsIn.Recipient]
			if asyncSmsReq == nil {
				log.Warning("No waiting for response request found for recipient %s", smsIn.Recipient)
			} else {
				delete(requests, smsIn.Recipient)
				asyncSmsReq.RespondWith(smsIn)
			}

		// Periodically check health
		case <-healthCheckTicker.C:
			log.Notice("Self diagnostics event. Num of Goroutines: %d", runtime.NumGoroutine())
			for _, asyncSmsReq := range requests {
				if asyncSmsReq.IsLonger(MAX_RESPONSE_WAITTIME) {
					log.Critical("Releasing request: %s", asyncSmsReq.SmsOut)
					delete(requests, asyncSmsReq.Originator)
					asyncSmsReq.Cancel()
				}
			}

		// print statistics
		case <-statsTicker.C:
			_totalProduced, _totalFailures =
				p.logStats(totalProduced, totalFailures, _totalProduced, _totalFailures)

		// registrations counter
		// if success is true, user successfully registered
		case success := <-p.countChan:
			totalProduced++
			if !success { totalFailures++ }
			if totalProduced == p.count {
				p.logSummary(totalProduced, totalFailures)
				close(done)
				return
			}

		// kill signal received
		case <-killChan:
			p.logSummary(totalProduced, totalFailures)
			<-p.semChan // clear blocking semaphore
			close(done)
			return
		}


	}

}

// Performs SMS request/response to remote Caipirinha server
// creates async sms requests and writes it to producer's async requests channel
//
// smsOut - message to send to server
// Returns read only channel to receive sms response from server
func (p *Producer) makeSmsRequest(smsOut *SmsOut) <-chan *SmsIn  {
	asyncReq, respChan := NewAsyncSmsReq(p.url, smsOut)
	p.asyncReqChan <- asyncReq
	return respChan
}

// TODO handle wanted response logic, when unexpected response from server can be received
func (p *Producer) makeRegistrationStep(phone, step string) bool {
	// make request and read response from response channel
	resp := <-p.makeSmsRequest(NewSmsOut(phone, step))
	if resp == nil {
		log.Error("Failed to perform registration step: %s, %s", phone, step)
		return false
	}
	return true
}

//
func (p *Producer) registerUser(phone, address string) {

	log.Debug("Registering and unregistering user with phone %s", phone)

	defer func() {
		<-p.semChan // release semaphore on exit
	}()

	// here is the multistepped registration process
	// if any step fails, the registration is considered failed
	if 	p.makeRegistrationStep(phone, REGISTER_CMD) &&
	   	p.makeRegistrationStep(phone, YES) &&
	   	p.makeRegistrationStep(phone, "25") &&
	   	p.makeRegistrationStep(phone, address) &&
	   	p.makeRegistrationStep(phone, "1") &&
	   	p.makeRegistrationStep(phone, UNREGISTER_CMD) {

		// report success
		log.Debug("Successfully registered and unregistered user with phone %s", phone)
		p.countChan <- true
	} else {
		// report failure
		log.Warning("Failed to register user with phone %s", phone)
		p.countChan <- false
	}

}

// Produces SMS user registrations
// killChan - channel to (gracefully) interrupt producer before desired number of registrations is made,
// or if the producer runs in infinite loop
func (p *Producer) Produce(killChan <-chan struct{}) {

	doneChan := make(chan struct{}) // write to this channel when desired number of registrations is made

	go p.handleRequests(doneChan, killChan) // startup request handler

	p.startTime = time.Now() // starting timer

	for i := 0; i != p.count; i++ {

		// not blocking channel reader
		// to interrupt producer loop earlier by Ctrl+C
		select {
			case <-doneChan: return
			default: // continue circle
		}

		p.semChan <- 1
		go p.registerUser(<-p.phoneGen, <-p.addressGen)
		runtime.Gosched() // Allow other goroutines to proceed
	}

	<-doneChan // wait until done
}