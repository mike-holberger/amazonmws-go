package amazonmwsapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"io"
	"net/url"
	"strings"
	"time"
)

// PrettPrintXML returns the XML data as a formatted string prepared for debugging output
func PrettPrintXML(data []byte) (string, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "  ")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			encoder.Flush()
			return string(b.Bytes()), nil
		}
		if err != nil {
			return "", err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return "", err
		}
	}
}

// ISO8601 format string for ISO8601
var ISO8601 = "2006-01-02T15:04:05Z"

// XMLTimestamp formats timestamp to ISO8601
func XMLTimestamp(t time.Time) string {
	return t.UTC().Format(ISO8601)
}

// CanonicalizedQueryString escapes querystring characters according to amzMWS documentation
func CanonicalizedQueryString(values url.Values) (str string) {
	// per aws docs and docs for values.Encode, we respect RFC 3986
	// we may not deal with utf-8, only ascii
	// params are sorted
	// we have to fix the '+' to '%20'
	str = values.Encode()
	str = strings.Replace(str, "+", "%20", -1)
	return
}

// Sign forms HMAC signature for authorizing requests to amzMWS
func Sign(str string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Region holds data about the target marketplace website
type Region struct {
	RegionID      string
	Country       string
	Endpoint      string
	MarketPlaceID string
}

// RegionByCountry returns Region data for a specified country
func RegionByCountry(country string) Region {
	for _, region := range Regions {
		if strings.EqualFold(region.Country, country) {
			return region
		}
	}
	panic("Invalid region, check your data")
}

// Regions is holding pre-loaded marketplace data
var Regions = []Region{
	{"NA", "US", "https://mws.amazonservices.com/", "ATVPDKIKX0DER"},
	{"NA", "CA", "https://mws.amazonservices.ca/", "A2EUQ1WTGCTBG2"},
	{"EU", "DE", "https://mws-eu.amazonservices.com/", "A1PA6795UKMFR9"},
	{"EU", "ES", "https://mws-eu.amazonservices.com/", "A1RKKUPIHCS9HS"},
	{"EU", "FR", "https://mws-eu.amazonservices.com/", "A13V1IB3VIYZZH"},
	{"EU", "IN", "https://mws.amazonservices.in/", "A21TJRUUN4KGV"},
	{"EU", "IT", "https://mws-eu.amazonservices.com/", "APJ6JRA9NG5V4"},
	{"EU", "UK", "https://mws-eu.amazonservices.com/", "A1F83G8C2ARO7P"},
	{"FE", "JP", "https://mws.amazonservices.jp/", "A1VC38T7YXB528"},
	{"CN", "CN", "https://mws.amazonservices.com.cn/", "AAHKV2X7AFYLW"},
}
