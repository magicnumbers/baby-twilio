Baby Twilio
===========

Most of the time with Twilio we want to do one of two things: send an SMS and
parse incoming SMSes. This package makes it easy to do both.


## Sending SMSes

In the most simple case, you can set three environment variables per Twilio
(`TWILIO_SID`, `TWILIO_AUTH_TOKEN`, `TWILIO_PHONE_NUMBER`) and then use a
little function to fire off your SMSes.

```go
package main

import (
	"baby-twilio"
	"log"
)

func main() {
	if _, err := twilio.SendSMS("+1-212-555-1212", "Well, hello there."); err != nil {
		log.Fatal("It didn‘t work: ", err)
	}
}
```

For more thorough examples, including ones that don‘t require setting
environment variables, see the `examples` directory.


## Receiving SMSes

Parsing incoming SMSes is just a matter of pointing a webhook at your server
in the Twilio admin interface and parsing incoming requests as such:

```go
package main

import (
	"github.com/meowgorithm/baby-twilio"
	"log"
	"net/http"
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

	http.ListenAndServe(":8000", nil)
}
```

For local development we recommend [ngrok][ng].


## License

MIT

[ng]: https://ngrok.com
