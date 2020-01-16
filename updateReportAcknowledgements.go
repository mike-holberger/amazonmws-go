package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"strconv"
)

// UpdateReportAcknowledgementsRequest acknowledges reports
type UpdateReportAcknowledgementsRequest struct {
	amazonRequest
}

// Acknowledged adds ack param to request
func (r *UpdateReportAcknowledgementsRequest) Acknowledged(ack bool) *UpdateReportAcknowledgementsRequest {
	r.params.Add("Acknowledged", strconv.FormatBool(ack))
	return r
}

// Do sends request to amazonMWS reports API
func (r *UpdateReportAcknowledgementsRequest) Do(ctx context.Context) (*UpdateReportAcknowledgementsResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &UpdateReportAcknowledgementsResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// UpdateReportAcknowledgementsResponse holds response data
type UpdateReportAcknowledgementsResponse struct {
	XMLName                            xml.Name `xml:"UpdateReportAcknowledgementsResponse"`
	Xmlns                              string   `xml:"xmlns,attr"`
	UpdateReportAcknowledgementsResult struct {
		Count      string       `xml:"Count"`
		ReportInfo []ReportInfo `xml:"ReportInfo"`
	} `xml:"UpdateReportAcknowledgementsResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}
