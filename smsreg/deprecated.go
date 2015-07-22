package smsreg
//
//import (
//	"github.com/op/go-logging"
//	"github.com/kolach/loadtest/smsin"
//	"github.com/kolach/loadtest/random"
//	"time"
//	"runtime"
//	"encoding/json"
//	"net/http"
//	"bytes"
//	"github.com/kolach/loadtest/postalcodes"
//	"os/signal"
//	"os"
//	"syscall"
//)
//
//const (
//	UNREGISTER_CMD 	= "SAIR"
//	REGISTER_CMD 	= "Register me"
//	APP_JSON 		= "application/json"
//	YES 			= "Si"
//
//	HEALTH_CHECK_INTERVAL = time.Second * 5
//	MAX_RESPONSE_WAITTIME = time.Second * 20
//	STATS_REPORT_INTERVAL = time.Second * 30
//)
//
//var log = logging.MustGetLogger("smsout")
//var startTime time.Time
//
//// Subscription
//// Used for sub/unsub to/from caipirinha server response to registration request
//type subscription struct {
//	Sms 			*SmsReq
//	RespChan 		chan<- *smsin.SmsIn 	// wait for response channel
//	SubscribedAt	time.Time				// time of subscription
//}
//
//// Request to send to caipirinha server
//type SmsReq struct {
//	Originator 	string `json:"originator"`
//	Receiver 	string `json:"receiver"`
//	Text 		string `json:"text"`
//}
//
//// Makes HTTP request
//// Returns response channel
//func makeSmsRequest(url string, originator, text string, subscribe chan<- *subscription) <-chan *smsin.SmsIn  {
//
//	respChan := make(chan *smsin.SmsIn)
//
//	smsReq  := &SmsReq{Originator: originator, Receiver: "9999999", Text: text}
//	sub 	:= &subscription{smsReq, respChan, time.Now()}
//
//	json, err := json.Marshal(smsReq)
//	if err != nil {
//		log.Error("Error marshaling json", err)
//		respChan <- nil
//	} else {
//		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
//		if err != nil {
//			log.Error("Error creating http request", err)
//			respChan <- nil
//		} else {
//			httpReq.Header.Set("Content-Type", APP_JSON)
//			httpReq.Header.Set("Accept", APP_JSON)
//			client := &http.Client{}
//			_, err = client.Do(httpReq)
//			if err != nil {
//				log.Error("Error making http request", err)
//				respChan <- nil
//			} else {
//				subscribe <- sub
//			}
//		}
//	}
//	return respChan
//}
//
//// Make HTTP request to SMS registration end point
////func makeSmsHttpRequest(url string, originator, text string) bool {
////
////	log.Debug("Making http request originator: %s, text: %s", originator, text)
////
////	smsReq := &SmsReq{Originator: originator, Receiver: "9999999", Text: text}
////	json, err := json.Marshal(smsReq)
////	if err != nil {
////		log.Error("Error marshaling json", err)
////		return false
////	}
////
////	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
////	if err != nil {
////		log.Error("Error creating http request", err)
////		return false
////	}
////
////
////	httpReq.Header.Set("Content-Type", APP_JSON)
////	httpReq.Header.Set("Accept", APP_JSON)
////
////	client := &http.Client{}
////	_, err = client.Do(httpReq)
////	if err != nil {
////		log.Error("Error making http request", err)
////		return false
////	}
////
////	return true
////
//////	defer resp.Body.Close()
//////	body, _ := ioutil.ReadAll(resp.Body)
//////	log.Debug("Response body is: %s", body)
////}
//
//// Make user registration HTTP request and block waitForResponseChan
//// waitForResponseChan will be unblocked on SMS message received for given phone number from caipirinha server
////func makeRegistrationStep(url string, respChan chan<- string, originator, text string) {
////	makeSmsHttpRequest(url, originator, text)
////	respChan <- text
////}
//
//
//// Multistep User Registration process
//// url - url to send request
//// phone - user's phone number
//// address - user's address
//// sem - semaphore channel
//// subscribe - subscribe to response channel
//// unsubscribe - unsubscribe from response channel
//func register(url, phone, address string, sem <-chan int, subscribe chan<- *subscription) {
//
//	// channel to block registration process untill the server responds
////	respChan := make(chan string)
//
//	// subscription to sub/unsub to/from caipirinha server response
////	subscription := NewSubscription(phone, respChan)
//
//	// 1. unsubscribe
//	// 2. close response channel
//	// 3. release semaphor
//	defer func() {
//		<-sem // release semaphore
//	}()
//
//	<-makeSmsRequest(url, phone, REGISTER_CMD, subscribe)
//	<-makeSmsRequest(url, phone, YES, subscribe)
//	<-makeSmsRequest(url, phone, "25", subscribe)
//	<-makeSmsRequest(url, phone, address, subscribe)
//	<-makeSmsRequest(url, phone, "1", subscribe)
//	<-makeSmsRequest(url, phone, UNREGISTER_CMD, subscribe)
//}
//
//
//// Registrations production circle
//func produce(url string, concurrency, count int, subscribe chan<- *subscription) {
//
//	log.Debug("Preparing system to produce registrations...")
//
//	// random generators
//	phoneGen, err := random.PhoneNumberGen("MX", "55")
//	if err != nil { panic(err) }
//	// it takes a while to load postal codes database
//	addressGen, err := random.AddressGen("MX")
//	if err != nil { panic(err) }
//
//	// semaphore channel. Limits number of concurrent users
//	sem := make(chan int, concurrency)
////	startTime = time.Now()
//
//	log.Debug("System is ready to produce user registrations")
//	startTime = time.Now()
//
//	for i := 0; i != count; i++ {
//		sem <- 1
//		go register(url, <-phoneGen, <-addressGen, sem, subscribe)
//		runtime.Gosched() // Allow other goroutines to proceed
//	}
//
//}
//
//func report(totalRegs, totalFails, prevTotalRegs, prevTotalFails int) (int, int) {
//	log.Notice("Time elapsed: %s, Users registered: %d (+%d), Failures: %d (+%d)",
//		time.Since(startTime), totalRegs, totalRegs - prevTotalRegs, totalFails, totalFails - prevTotalFails)
//	return totalRegs, totalFails
//}
//
//// Consumption circle
//func consume(resp <-chan *smsin.SmsIn, subscribe <-chan *subscription, count int, done chan<- bool) {
//
//	totalRegs  		:= 0 // total users registered
//	totalFails 		:= 0 // total users failed to register
//	prevTotalRegs  	:= 0 // total users registered
//	prevTotalFails 	:= 0 // total users failed to register
//
//	healthCheckTicker := time.NewTicker(HEALTH_CHECK_INTERVAL)
//	statsReportTicker := time.NewTicker(STATS_REPORT_INTERVAL)
//
//	signalChan := make(chan os.Signal,  1)	// chanel to send kill signal
//	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM) // handle Ctrl^C and SIGTERM
//
//	subscriptions := make(map[string] *subscription)
//
//	for {
//		select {
//
//		// subscribe to server response
//		case s := <-subscribe:
//			log.Debug("Subscribing channel for recepient %s", s.Sms.Originator)
//			subscriptions[s.Sms.Originator] = s
//
//		// unsubscribe from server response
////		case s := <-unsubscribe:
////			log.Debug("Unsubscribing channel for recepient %s", s.Recipient)
////			delete(subscriptions, s.Recipient)
////			totalRegs++
////			if count == totalRegs {
////				report(totalRegs, totalFails, prevTotalRegs, prevTotalFails)
////				done <- true
////			}
//
//
//		// server response: find subscription and read channel,
//		// releasing it to make a next registration step
//		case r := <-resp:
//			s := subscriptions[r.Recipient]
//			if s == nil {
//				log.Warning("No subscription found for msg %s", r)
//			} else {
//				// response received, stop blocking next reg step
//				s.RespChan <- r
//			}
//
//		// print statistics
//		case <-statsReportTicker.C:
//			prevTotalRegs, prevTotalFails = report(totalRegs, totalFails, prevTotalRegs, prevTotalFails)
//
//		// Health check
//		// Check subscriptions waiting for server response and close them in case
//		// the wait time is more than MAX_RESPONSE_WAITTIME
//		case <-healthCheckTicker.C:
//			log.Notice("Self diagnostics event. Num goroutines: %d", runtime.NumGoroutine())
//			for _, s := range subscriptions {
//				if time.Since(s.SubscribedAt) > MAX_RESPONSE_WAITTIME {
//					totalFails++    // increment failures counter and report the problem
//					log.Critical("Releasing request for recepient: %s, msg: %s", s.Sms.Originator, s.Sms.Text)
//					s.RespChan <- nil
//				}
//			}
//
//
//		case <-signalChan:
//			report(totalRegs, totalFails, prevTotalRegs, prevTotalFails)
//			done <- true
//		}
//
//	}
//}
//
//// resp - channel to listen for server responses
//func Produce(url string, concurrency, count int, resp <-chan *smsin.SmsIn) {
//
//	postalcodes.GetDb("MX") // preload mexican db
//
//	log.Notice("Producing SMS user registration with max number of concurrent users: %d", concurrency)
//	log.Notice("Total user registrations to produce: %d", count)
//	log.Notice("Uri: %s", url)
//
//	subscribe 	:= make(chan *subscription, concurrency) // subscribe channel
//
//	done  		:= make(chan bool)
//
//	go consume(resp, subscribe, count, done)
//	go produce(url, concurrency, count, subscribe)
//
//	<-done
//}
//
