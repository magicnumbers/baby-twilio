package twilio

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL        = "https://api.twilio.com/2010-04-01/Accounts/"
	suffix         = "/Messages.json"
	defaultTimeout = time.Second * 10
)

// NewClient is a factory for the client we use to send requests. One could
// also just initialize the `Client` struct manually too, of course,
func NewClient(sid, authToken, phoneNumber string) *Client {
	return &Client{
		SID:         sid,
		AuthToken:   authToken,
		PhoneNumber: phoneNumber,
	}
}

// Client manages our requests
type Client struct {
	SID         string
	AuthToken   string
	PhoneNumber string
	Timeout     time.Duration
}

// Find out if required values on the struct have been set
func (c Client) checkSettings() error {
	var e []string
	if c.SID == "" {
		e = append(e, "SID")
	}
	if c.AuthToken == "" {
		e = append(e, "auth token")
	}
	if c.PhoneNumber == "" {
		e = append(e, "phone number")
	}

	if len(e) > 0 {
		// Make nice error string (no Oxford comma!)
		errs := ""
		for i := 0; i < len(e); i++ {
			if i == 0 {
				errs += e[i]
			} else if i == len(e)-1 {
				errs += " and " + e[i]
			} else {
				errs += ", " + e[i]
			}
		}
		errText := "some things aren't set: " + errs
		return errors.New(errText)
	}

	return nil
}

// SendSMS sends an SMS to a phone number
func (c Client) SendSMS(to string, msg string) (twilioResponse Response, err error) {
	if err = c.checkSettings(); err != nil {
		return twilioResponse, err
	}

	// Set default request timeout if one hasn't been specified
	httpClient := &http.Client{}
	if c.Timeout == 0 {
		httpClient.Timeout = defaultTimeout
	}

	// Set form values
	v := url.Values{}
	v.Set("To", to)
	v.Set("From", c.PhoneNumber)
	v.Set("Body", msg)
	rb := *strings.NewReader(v.Encode())

	// API url
	url := baseURL + c.SID + suffix

	// Prepare the request
	req, err := http.NewRequest("POST", url, &rb)
	if err != nil {
		return twilioResponse, err
	}
	req.SetBasicAuth(c.SID, c.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Fire off the request
	res, err := httpClient.Do(req)

	// Parse response
	if err = json.NewDecoder(res.Body).Decode(&twilioResponse); err != nil {
		return twilioResponse, err
	}

	return twilioResponse, twilioResponse.CheckError()
}
