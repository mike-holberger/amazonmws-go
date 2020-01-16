package amazonmwsapi

import (
	"context"
	"encoding/xml"
)

// RequestReportRequest requests a single amzMWS report be prepared for download
type RequestReportRequest struct {
	amazonRequest
}

// Do sends request to amazonMWS reports API and returns report request info
func (r *RequestReportRequest) Do(ctx context.Context) (*RequestReportResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &RequestReportResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// RequestReportResponse type
type RequestReportResponse struct {
	XMLName             xml.Name `xml:"RequestReportResponse"`
	Xmlns               string   `xml:"xmlns,attr"`
	RequestReportResult struct {
		ReportRequestInfo struct {
			ReportType             string `xml:"ReportType"`
			ReportProcessingStatus string `xml:"ReportProcessingStatus"`
			EndDate                string `xml:"EndDate"`
			Scheduled              string `xml:"Scheduled"`
			ReportRequestID        string `xml:"ReportRequestId"`
			SubmittedDate          string `xml:"SubmittedDate"`
			StartDate              string `xml:"StartDate"`
		} `xml:"ReportRequestInfo"`
	} `xml:"RequestReportResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}
