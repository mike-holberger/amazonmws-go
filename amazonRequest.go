package amazonmwsapi

import (
	"bytes"
	"net/url"
)

type amazonRequest struct {
	endpoint string
	params   url.Values
	method   string
	body     *bytes.Buffer
	client   *AmazonClient
}
