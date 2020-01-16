package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"fmt"
)

// GetReportRequestListRequest calls a list of requested reports and thier status
type GetReportRequestListRequest struct {
	amazonRequest
}

// ReportRequestIDList adds list of report IDs to request - not required
func (r *GetReportRequestListRequest) ReportRequestIDList(reportIDs []string) *GetReportRequestListRequest {
	for i, id := range reportIDs {
		r.params[fmt.Sprintf("ReportRequestIdList.Id.%d", (i+1))] = []string{id}
	}
	return r
}

// ReportTypeList adds list of reportTypes to request - not required
func (r *GetReportRequestListRequest) ReportTypeList(reportTypes []string) *GetReportRequestListRequest {
	for i, t := range reportTypes {
		r.params[fmt.Sprintf("ReportTypeList.Type.%d", (i+1))] = []string{t}
	}
	return r
}

// Do sends request to amazonMWS reports API and returns report request info
func (r *GetReportRequestListRequest) Do(ctx context.Context) (*GetReportRequestListResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &GetReportRequestListResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// GetReportRequestListResponse type
type GetReportRequestListResponse struct {
	XMLName                    xml.Name `xml:"GetReportRequestListResponse"`
	Xmlns                      string   `xml:"xmlns,attr"`
	GetReportRequestListResult struct {
		NextToken         string `xml:"NextToken"`
		HasNext           string `xml:"HasNext"`
		ReportRequestInfo []struct {
			ReportRequestID        string `xml:"ReportRequestId"`
			ReportType             string `xml:"ReportType"`
			StartDate              string `xml:"StartDate"`
			EndDate                string `xml:"EndDate"`
			Scheduled              string `xml:"Scheduled"`
			SubmittedDate          string `xml:"SubmittedDate"`
			ReportProcessingStatus string `xml:"ReportProcessingStatus"`
			GeneratedReportID      string `xml:"GeneratedReportId"`
			StartedProcessingDate  string `xml:"StartedProcessingDate"`
			CompletedDate          string `xml:"CompletedDate"`
		} `xml:"ReportRequestInfo"`
	} `xml:"GetReportRequestListResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}
