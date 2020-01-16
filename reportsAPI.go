package amazonmwsapi

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

const reportsAPIversion = "2009-01-01"

// ReportsAPI hold calls for downloading (order) reports
type ReportsAPI struct {
	client   *AmazonClient
	endpoint string
}

// NewReportsAPI creates and configures new ReportsAPI object
func NewReportsAPI(cli *AmazonClient) *ReportsAPI {
	api := &ReportsAPI{client: cli}
	api.endpoint = api.client.Region.Endpoint + "Reports/" + reportsAPIversion
	return api
}

// GetReportList calls a list of reports available for download
func (api *ReportsAPI) GetReportList() *GetReportListRequest {
	return &GetReportListRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params:   url.Values{"Action": {"GetReportList"}, "Version": {reportsAPIversion}},
			method:   "POST",
		},
	}
}

// GetReport downloads a single report
func (api *ReportsAPI) GetReport(reportID string) *GetReportRequest {
	return &GetReportRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params: url.Values{"Action": {"GetReport"}, "Version": {reportsAPIversion},
				"ReportId": {reportID}},
			method: "POST",
		},
	}
}

// UpdateReportAcknowledgements downloads a single report
func (api *ReportsAPI) UpdateReportAcknowledgements(reportIDs []string) *UpdateReportAcknowledgementsRequest {
	req := &UpdateReportAcknowledgementsRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params:   url.Values{"Action": {"UpdateReportAcknowledgements"}, "Version": {reportsAPIversion}},
			method:   "POST",
		},
	}
	for i, rep := range reportIDs {
		key := fmt.Sprintf("ReportIdList.Type.%d", (i + 1))
		req.params.Add(key, rep)
	}
	return req
}

// RequestReport requests a report be prepared by amazonMWS
func (api *ReportsAPI) RequestReport(reportType string) *RequestReportRequest {
	return &RequestReportRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params: url.Values{"Action": {"RequestReport"}, "Version": {reportsAPIversion},
				"ReportType": {reportType}},
			method: "POST",
		},
	}
}

// GetReportRequestList calls a list of requested reports and thier status
func (api *ReportsAPI) GetReportRequestList() *GetReportRequestListRequest {
	return &GetReportRequestListRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params:   url.Values{"Action": {"GetReportRequestList"}, "Version": {reportsAPIversion}},
			method:   "POST",
		},
	}
}

// DownloadInvReport requests an inventory report, checks report progress and downloads report to destination
func (api *ReportsAPI) DownloadInvReport(ctx context.Context, reportType string, filePath string) error {
	// request product report
	reqRepResp, err := api.RequestReport(reportType).Do(ctx)
	if err != nil {
		return err
	}

	// get report request list
	reportReady := false
	var genRepID string

	for !reportReady {
		repReqListResp, err := api.GetReportRequestList().
			ReportRequestIDList([]string{reqRepResp.RequestReportResult.ReportRequestInfo.ReportRequestID}).
			Do(ctx)
		if err != nil {
			return err
		}

		if len(repReqListResp.GetReportRequestListResult.ReportRequestInfo) == 0 ||
			repReqListResp.GetReportRequestListResult.ReportRequestInfo[0].ReportProcessingStatus != "_DONE_" {
			// Report not ready, try again later (strict request limit here => 15 calls)
			time.Sleep(4 * time.Second)
			continue
		}

		reportReady = true
		genRepID = repReqListResp.GetReportRequestListResult.ReportRequestInfo[0].GeneratedReportID
	}

	// Download generated report
	err = api.GetReport(genRepID).Download(ctx, filePath)
	if err != nil {
		return err
	}

	return nil
}
