package twilio

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// Request represents an request to a Twilio webhook. These seem to always
// come in as a standard set of key/value form pairs (as opposed to JSON).
type Request struct {
	ToCity        string `url:"ToCity"`
	FromZip       int    `url:"FromZip"`
	MessageSID    string `url:"MessageSID"`
	SMSMessageSID string `url:"SMSMessageSID"`
	ToCountry     string `url:"ToCountry"`
	FromCountry   string `url:"FromCountry"`
	ToZip         string `url:"ToZip"`
	FromState     string `url:"FromState"`
	SMSStatus     string `url:"SMSStatus"` // should maybe be an enum?
	FromCity      string `url:"FromCity"`
	Body          string `url:"Body"`
	To            string `url:"To"`
	AccountSID    string `url:"AccountSID"`
	NumMedia      int    `url:"NumMedia"`
	From          string `url:"From"`
	SMSSid        string `url:"SMSSid"`
	NumSegments   int    `url:"NumSegments"`
}

// NewDecoder is simply a wrapper around `Request.Decode()`. Use it like so:
//
//     req, _ := twilio.NewDecoder(htmlRequest)
//     fmt.Println("Hey look, this message is from", req.From)
//
func NewDecoder(r *http.Request) (*Request, error) {

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Parse query string data
	var twilioReq Request
	err = twilioReq.Decode(body)
	return &twilioReq, err
}

// Decode decodes our struct
// TODO: Fix cyclomatic complexity
func (r *Request) Decode(queryParams []byte) error {

	const tagName = "url"

	// Parse query params
	params, err := url.ParseQuery(string(queryParams))
	if err != nil {
		return err
	}

	t := reflect.TypeOf(*r)

	// Iterate over all fields in the struct
	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		// For debugging
		//fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)

		// `field.Name` is the actual field name
		// `field.Type.Name()` is the tag type (i.e. "string")
		// `tag` is the name of the tag

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		// NOTE: Values in the query params map are stored in a slice, usually
		// of a length of exactly 1, so we have to also get that first value.
		// For example, to get the "From" field we'd normally access it like:
		//
		//     from := params["From"][0]
		//
		if val, ok := params[tag]; ok {

			// Skip if this isn't a slice for some reason
			if len(val) < 1 {
				continue
			}

			// Now we make sure the param is a slice of at least one element
			v := reflect.ValueOf(val)
			switch v.Kind() {
			case reflect.Slice:

				structVal := reflect.ValueOf(r).Elem()
				structFieldVal := structVal.FieldByName(field.Name)

				if !structFieldVal.IsValid() {
					return fmt.Errorf("no such field: %s in obj", tag)
				}

				if !structFieldVal.CanSet() {
					return fmt.Errorf("cannot set %s field value", tag)
				}

				switch field.Type.Name() {
				case "string":
					structFieldVal.Set(reflect.ValueOf(val[0]))
					break
				case "int":
					i, err := strconv.Atoi(val[0])
					if err != nil {
						return err
					}
					structFieldVal.Set(reflect.ValueOf(i))
					break
				default:
					log.Printf("warning, no serializer set for %v", field.Type.Name())
				}

			}

		}

	}

	return nil
}
