package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
)

// GetReportListRequest gets MWS reports available for download
type GetReportListRequest struct {
	amazonRequest
}

// ReportTypes add requested report types to request
func (r *GetReportListRequest) ReportTypes(reportTypes []string) *GetReportListRequest {
	for i, rep := range reportTypes {
		key := fmt.Sprintf("ReportTypeList.Type.%d", (i + 1))
		r.params.Add(key, rep)
	}
	return r
}

// Acknowledged adds ack param to request
func (r *GetReportListRequest) Acknowledged(ack bool) *GetReportListRequest {
	r.params.Add("Acknowledged", strconv.FormatBool(ack))
	return r
}

// Do sends request to amazonMWS reports API
func (r *GetReportListRequest) Do(ctx context.Context) (*GetReportListResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &GetReportListResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoNext sends gets next page of amazon reports if HasNext == true
func (r *GetReportListRequest) DoNext(ctx context.Context, nextToken string) (*GetReportListByNextTokenResponse, error) {
	nextReq := &amazonRequest{
		client:   r.client,
		endpoint: r.endpoint,
		method:   r.method,
		params: url.Values{
			"Action":    {"GetReportListByNextToken"},
			"Version":   r.params["Version"],
			"NextToken": {nextToken},
		},
	}

	respBytes, err := r.client.callAPI(ctx, nextReq)
	if err != nil {
		return nil, err
	}

	xmlResponse := &GetReportListByNextTokenResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoAll calls multiple page response if necessary
func (r *GetReportListRequest) DoAll(ctx context.Context) ([]ReportInfo, error) {
	reports := []ReportInfo{}
	firstResp, err := r.Do(ctx)
	if err != nil {
		return nil, err
	}
	reports = append(reports, firstResp.GetReportListResult.ReportInfo...)

	hasNext := firstResp.GetReportListResult.HasNext
	nextToken := firstResp.GetReportListResult.NextToken
	for hasNext {
		nextResp, err := r.DoNext(ctx, nextToken)
		if err != nil {
			return reports, err
		}
		reports = append(reports, nextResp.GetReportListByNextTokenResult.ReportInfo...)

		nextToken = nextResp.GetReportListByNextTokenResult.NextToken
		hasNext = nextResp.GetReportListByNextTokenResult.HasNext
	}

	return reports, nil
}

// ReportInfo contains report info
type ReportInfo struct {
	ReportID        string `xml:"ReportId"`
	ReportType      string `xml:"ReportType"`
	ReportRequestID string `xml:"ReportRequestId"`
	AvailableDate   string `xml:"AvailableDate"`
	Acknowledged    string `xml:"Acknowledged"`
}

// GetReportListResponse holds response data
type GetReportListResponse struct {
	XMLName             xml.Name `xml:"GetReportListResponse"`
	Xmlns               string   `xml:"xmlns,attr"`
	GetReportListResult struct {
		NextToken  string       `xml:"NextToken"`
		HasNext    bool         `xml:"HasNext"`
		ReportInfo []ReportInfo `xml:"ReportInfo"`
	} `xml:"GetReportListResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// GetReportListByNextTokenResponse holds response data
type GetReportListByNextTokenResponse struct {
	XMLName                        xml.Name `xml:"GetReportListByNextTokenResponse"`
	Xmlns                          string   `xml:"xmlns,attr"`
	GetReportListByNextTokenResult struct {
		NextToken  string       `xml:"NextToken"`
		HasNext    bool         `xml:"HasNext"`
		ReportInfo []ReportInfo `xml:"ReportInfo"`
	} `xml:"GetReportListByNextTokenResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}
