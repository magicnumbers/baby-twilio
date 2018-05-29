package twilio

import "os"

// Environment variable keys
const (
	TwilioSID         = "TWILIO_SID"
	TwilioAuthToken   = "TWILIO_AUTH_TOKEN"
	TwilioPhoneNumber = "TWILIO_PHONE_NUMBER"
)

// SendSMS sends an SMS. We expect the following environment varaibles to be
// set: TWILIO_SID, TWILIO_AUTH_TOKEN, TWILIO_PHONE_NUMBER.
//
// To send an SMS without using environment variables, see `Client`.
func SendSMS(to string, msg string) (Response, error) {
	return Client{
		SID:         os.Getenv(TwilioSID),
		AuthToken:   os.Getenv(TwilioAuthToken),
		PhoneNumber: os.Getenv(TwilioPhoneNumber),
	}.SendSMS(to, msg)
}
