package main

import (
	"baby-twilio"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var (
		to  string
		err error
		res twilio.Response
	)

	// Ask for a phone number
	fmt.Printf("What‘s ur phone number? ")
	r := bufio.NewReader(os.Stdin)
	if to, err = r.ReadString('\n'); err != nil {
		fmt.Println("Something went wrong: ", err)
		os.Exit(1)
	}

	client := twilio.Client{
		SID:         os.Getenv("TWILIO_SID"),
		AuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
		PhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
	}

	res, err = client.SendSMS(to, "Oh, hello.")

	// Print out the response regardless of weather or not there was an error
	if b, err := json.MarshalIndent(res, "", "  "); err == nil {
		fmt.Printf("\nTwilio says:\n\n")
		fmt.Printf("%s\n\n", string(b))
	} else {
		fmt.Println("Couldn‘t Marshal response into JSON ", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("It didn‘t work: ", err)
		os.Exit(1)
	}

	fmt.Println("SMS sent successfully.")

}
