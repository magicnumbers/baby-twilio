Baby Twilio
===========

Most of the time with Twilio we want to do one of two things: send an SMS and
parse incoming SMSes. This package makes it easy to to both.

But hold on a sec: this is still in development.


## Simple Example

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
