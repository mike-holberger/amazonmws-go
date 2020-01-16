package amazonmwsapi

import (
	"bytes"
	"encoding/xml"
)

// ErrorResponse holds error data from errored API call
type ErrorResponse struct {
	XMLName xml.Name `xml:"ErrorResponse"`
	Xmlns   string   `xml:"xmlns,attr"`
	Errors  []struct {
		Type    string `xml:"Type"`
		Code    string `xml:"Code"`
		Message string `xml:"Message"`
	} `xml:"Error"`
	RequestID string `xml:"RequestId"`
}

func (e *ErrorResponse) Error() string {
	buf := bytes.Buffer{}
	buf.WriteString("ERRORS from Amazon API: ")

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "; ")
	}

	return buf.String()
}
