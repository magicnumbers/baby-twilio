package main

/*

	Example of reading incoming SMSes sent to a webhook as defined in the
	Twilio admin interface.

*/

import (
	"log"
	"net/http"

	twilio "github.com/meowgorithm/baby-twilio"
)

func main() {

	// Handle incoming HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			sms *twilio.Request
			err error
		)

		// Decode request data from Twilio
		if sms, err = twilio.NewDecoder(r); err != nil {
			log.Println("could not decode incoming SMS: ", err)
			return
		}

		log.Printf("Incoming SMS from %s: %s", sms.From, sms.Body)
	})

	// Run webserver
	http.ListenAndServe(":8000", nil)
}
