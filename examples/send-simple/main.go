package main

/*

	This example assumes that a few environment variables are set:

	* TWILIO_SID
	* TWILIO_AUTH_TOKEN
	* TWILIO_PHONE_NUMBER

*/

import (
	"bufio"
	"fmt"
	"os"

	twilio "github.com/meowgorithm/baby-twilio"
)

func main() {
	var (
		to  string
		err error
	)

	// Ask for a phone number
	fmt.Printf("What‘s ur phone number? ")
	r := bufio.NewReader(os.Stdin)
	if to, err = r.ReadString('\n'); err != nil {
		fmt.Println("Something went wrong: ", err)
		os.Exit(1)
	}

	// Send SMS
	if _, err := twilio.SendSMS(to, "Well, hello there."); err != nil {
		fmt.Println("It didn‘t work: ", err)
		os.Exit(1)
	}

	fmt.Println("Your SMS was sent!")
}
