package amazonmwsapi

import "net/url"

const feedsAPIversion = "2009-01-01"

// FeedsAPI hold calls for downloading (order) reports
type FeedsAPI struct {
	client   *AmazonClient
	endpoint string
}

// NewFeedsAPI creates and configures new ReportsAPI object
func NewFeedsAPI(cli *AmazonClient) *FeedsAPI {
	api := &FeedsAPI{client: cli}
	api.endpoint = api.client.Region.Endpoint + "Feeds/" + feedsAPIversion
	return api
}

// SubmitFeed posts feed content in XML to amazonMWS
func (api *FeedsAPI) SubmitFeed() *SubmitFeedRequest {
	return &SubmitFeedRequest{
		amazonRequest: amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params:   url.Values{"Action": {"SubmitFeed"}, "Version": {feedsAPIversion}},
			method:   "POST",
		},
	}
}
