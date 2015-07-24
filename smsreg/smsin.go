package smsreg


import "fmt"

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

