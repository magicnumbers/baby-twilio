package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// This is what Twilio's time strings look like
	twilioTimeLayout = "Mon, 2 Jan 2006 15:04:05 -0700"
)

// Response contains values in a Twilio success response received after sending
// an SMS.
type Response struct {

	// Error indicates weather or not the response was an error. This is not
	// part of a normal Twilio response: we added it ourselves to remove some
	// ambigulity in our repsonses
	Error bool `json:"error"`

	// Found only in error responses
	Code     int    `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	MoreInfo string `json:"more_info,omitempty"`

	// This is an ambigulous response. If the response is an error, this is an
	// `int`, otherwise this is a `string` (often "queued"). To deal with this
	// curious API design we don't export this and instead process the response
	// into either `.HTTPStatusCode` for errors or `.MessageStatus` for
	// successes.
	Status interface{} `json:"status,omitempty"`

	// HTTPStatusCode is the JSON `status` key in error responses. It will not
	// be set in success responses.
	HTTPStatusCode int `json:"http_status_code,omitempty"`

	// MessageStatus is the JSON `status` key in success responses. It will not
	// be set in error responses.
	MessageStatus string `json:"message_status,omitempty"`

	// Time strings we need to parse separately
	DateCreated string `json:"date_created,omitempty"`
	DateUpdated string `json:"date_updated,omitempty"`
	DateSent    string `json:"date_sent,omitempty"`

	// Parsed time strings
	Created *time.Time `json:"created,omitempty"`
	Updated *time.Time `json:"updated,omitempty"`
	Sent    *time.Time `json:"sent,omitempty"`

	// Found in success responses
	SID                 string            `json:"sid,omitempty"`
	AccountSID          string            `json:"account_sid,omitempty"`
	To                  string            `json:"to,omitempty"`
	From                string            `json:"from,omitempty"`
	MessagingServiceSID string            `json:"messaging_service_sid,omitempty"`
	Body                string            `json:"body,omitempty"`
	NumSegments         string            `json:"num_segments,omitempty"` // Twilio, why is this a string?
	NumMedia            string            `json:"num_media,omitempty"`    // Twilio, why is this a string?
	Direction           string            `json:"direction,omitempty"`
	APIVersion          string            `json:"api_version,omitempty"`
	Price               string            `json:"price,omitempty"`
	PriceUnit           string            `json:"price_unit,omitempty"`
	ErrorCode           string            `json:"error_code,omitempty"`
	ErrorMesssage       string            `json:"error_messsage,omitempty"`
	URI                 string            `json:"uri,omitempty"`
	SubResourceURIs     map[string]string `json:"sub_resource_uris,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshalling for this struct. We do it
// here to handle the ambiguity of certain fields.
func (r *Response) UnmarshalJSON(b []byte) (err error) {
	// To avoid a recursive unmarhsalling loop and stack overflow we make a
	// copy of our struct and unmarshal into that
	type Alias Response
	alias := &struct {
		*Alias
	}{
		(*Alias)(r), // copy our struct
	}

	if err = json.Unmarshal(b, &alias); err != nil {
		return err
	}

	// Did we recieve an error response?
	if alias.Code != 0 {
		r.Error = true
	}

	// Figure out what 'status' is supposed to be
	switch v := alias.Status.(type) {
	case float64: // why is this not an int?
		r.HTTPStatusCode = int(v)
	case string:
		r.MessageStatus = v
	}

	// Parse times
	if r.Created, err = parseTime(alias.DateCreated); err != nil {
		return err
	}
	if r.Updated, err = parseTime(alias.DateUpdated); err != nil {
		return err
	}
	if r.Sent, err = parseTime(alias.DateSent); err != nil {
		return err
	}

	// Unset the string values
	r.DateCreated = ""
	r.DateUpdated = ""
	r.DateSent = ""

	return nil
}

// CheckError checks the response object for an error. If we recieved one,
// return an error
func (r Response) CheckError() error {
	if r.Error {
		return fmt.Errorf("%d %s: error %d (for more info see %s)",
			r.HTTPStatusCode,
			http.StatusText(r.HTTPStatusCode),
			r.Code,
			r.MoreInfo)
	}
	return nil
}

// Parse Twilio time strings into native Time types
func parseTime(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}
	t, err := time.Parse(twilioTimeLayout, str)
	if err != nil {
		return nil, err
	}
	return &t, err
}
