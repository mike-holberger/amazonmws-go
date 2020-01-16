package amazonmwsapi

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

// Creds holds amzMWS client credential data
type Creds struct {
	AccessID    string
	AccessKey   string
	CompanyName string
	Merchant    string
}

// AmazonClient executes requests to the amzMWS api
type AmazonClient struct {
	credentials      Creds
	Region           Region
	SignatureMethod  string
	SignatureVersion string
	UserAgent        string
	Logger           *logrus.Logger
}

// NewAmazonClient creates and configures AmazonClient
func NewAmazonClient(creds Creds, countryCode string, log *logrus.Logger) *AmazonClient {
	h, _ := os.Hostname()
	return &AmazonClient{
		credentials:      creds,
		SignatureVersion: "2",
		SignatureMethod:  "HmacSHA256",
		Region:           RegionByCountry(countryCode),
		UserAgent:        fmt.Sprintf("%s/amazonAlert (Language=go; Host=%s)", creds.CompanyName, h),
		Logger:           log,
	}
}

func (c *AmazonClient) parseRequest(req *amazonRequest) (*http.Request, error) {
	if c.credentials.AccessID == "" || c.credentials.AccessKey == "" || c.credentials.Merchant == "" || c.Region.Endpoint == "" {
		err := errors.New("Incomplete Request")
		return nil, err
	}

	req.params.Add("SellerId", c.credentials.Merchant)
	req.params.Add("AWSAccessKeyId", c.credentials.AccessID)
	req.params.Add("MarketplaceId.Id.1", c.Region.MarketPlaceID)
	req.params.Add("SignatureMethod", c.SignatureMethod)
	req.params.Add("SignatureVersion", c.SignatureVersion)
	req.params.Add("Timestamp", XMLTimestamp(time.Now()))

	stringToSign, err := c.stringToSign(req)
	if err != nil {
		return nil, err
	}
	signature := Sign(stringToSign, []byte(c.credentials.AccessKey))
	req.params.Add("Signature", signature)

	url, err := url.Parse(req.endpoint)
	if err != nil {
		return nil, err
	}
	url.RawQuery = CanonicalizedQueryString(req.params)

	request, err := http.NewRequest(req.method, url.String(), nil)
	if req.body != nil {
		request, err = http.NewRequest(req.method, url.String(), req.body)
		request.Header.Add("Content-Type", "text/xml")
	}

	request.Header.Add("User-Agent", c.UserAgent)
	return request, err
}

func (c *AmazonClient) callAPI(ctx context.Context, req *amazonRequest) ([]byte, error) {
	// Parse request params
	request, err := c.parseRequest(req)
	if err != nil {
		return nil, err
	}

	if c.Logger != nil {
		fmt.Println("REQUESTING: ", request.URL.String())
	}

	// Send request to amzMWS api
	request.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		if c.Logger != nil {
			c.Logger.WithFields(structs.Map(err)).WithField("Requested", request.URL.String()).
				Error("FAILED Amazon callAPI: " + err.Error())
		}
		return nil, err
	}

	// Read response
	bodyContents, _ := ioutil.ReadAll(resp.Body)
	defer func() {
		cerr := resp.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	// Check for http error
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(bodyContents))
	}

	if c.Logger != nil {
		c.Logger.WithField("Requested", request.URL.String()).
			Debug("RESPONSE: ", string(bodyContents))
	}

	return bodyContents, nil
}

func (c *AmazonClient) downloadReport(ctx context.Context, req *amazonRequest, filepath string) error {
	// Parse request params
	request, err := c.parseRequest(req)
	if err != nil {
		return err
	}

	if c.Logger != nil {
		c.Logger.Debug("REQUESTING: ", request.URL.String())
	}

	// Send request to amzMWS api
	request.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		if c.Logger != nil {
			c.Logger.WithFields(structs.Map(err)).WithField("Requested", request.URL.String()).
				Error("FAILED Amazon callAPI: " + err.Error())
		}
		return err
	}

	// Check for http error
	if resp.StatusCode != 200 {
		// Read and return http error response
		bodyContents, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf(string(bodyContents))
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file simultaneous with download stream
	_, err = io.Copy(out, resp.Body)
	return err
}

func (c *AmazonClient) submitFeed(ctx context.Context, r *SubmitFeedRequest) ([]byte, error) {
	// Encode feed into xml and attach body
	teeBytes := new(bytes.Buffer)
	teeBytes.Write([]byte(xml.Header))
	err := xml.NewEncoder(teeBytes).Encode(r.feed)
	if err != nil {
		if c.Logger != nil {
			c.Logger.WithFields(structs.Map(err)).WithFields(structs.Map(r.feed)).
				Error("FAILED xml encode during Amazon submitFeed: " + err.Error())
		}
		return nil, err
	}

	// Use teeReader to read bytes for MD5, and simultaneously copy into request body
	r.body = new(bytes.Buffer)
	tee := io.TeeReader(teeBytes, r.body)

	// Calculate and attach MD5 sum of XML body
	hash := md5.New()
	if _, err := io.Copy(hash, tee); err != nil {
		if c.Logger != nil {
			c.Logger.WithFields(structs.Map(err)).WithFields(structs.Map(r.feed)).
				Error("FAILED to copy into MD5 hash during Amazon submitFeed: " + err.Error())
		}
		return nil, err
	}
	contMD5base64 := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	r.params["ContentMD5Value"] = []string{contMD5base64}

	//fmt.Printf("\n\n[[[[%s]]]]\n\n", r.body.String())

	// Submit feed
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		// callAPI() logs client errors
		return nil, err
	}

	return respBytes, err
}

func (c *AmazonClient) parseAPIerrors(resp []byte) error {
	errResponse := &ErrorResponse{}
	err := xml.Unmarshal(resp, errResponse)
	if err != nil {
		return fmt.Errorf("UNABLE TO UNMARSHAL API RESPONSE: %s", err.Error())
	}

	return errResponse
}

func (c *AmazonClient) stringToSign(req *amazonRequest) (stringToSign string, err error) {
	endpoint, err := url.Parse(req.endpoint)
	if err != nil {
		return
	}
	stringToSign = strings.Join([]string{
		req.method,
		strings.ToLower(endpoint.Host),
		endpoint.Path,
		CanonicalizedQueryString(req.params),
	}, "\n")

	return
}
