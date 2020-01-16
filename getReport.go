package amazonmwsapi

import (
	"bytes"
	"context"
	"encoding/csv"
)

// GetReportRequest requests a single amzMWS report for download
type GetReportRequest struct {
	amazonRequest
}

// Do sends request to amazonMWS reports API and returns report data maps
func (r *GetReportRequest) Do(ctx context.Context) ([]map[string]string, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	return r.parseTSVData(respBytes)
}

func (r *GetReportRequest) parseTSVData(rep []byte) ([]map[string]string, error) {
	// Parse TSV into string data slices
	tsvReader := csv.NewReader(bytes.NewReader(rep))
	tsvReader.Comma = '\t'
	tsvReader.FieldsPerRecord = -1
	data, err := tsvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Parse into maps indexed by column header title
	rowMaps := make([]map[string]string, (len(data) - 1))
	for rowIndex, dataRow := range data[1:] {
		rowMap := make(map[string]string)
		for colIndex, header := range data[0] {
			rowMap[header] = dataRow[colIndex]
		}
		rowMaps[rowIndex] = rowMap
	}
	return rowMaps, nil
}

// Download sends request to amazonMWS reports API and downloads report to filepath
func (r *GetReportRequest) Download(ctx context.Context, filePath string) error {
	err := r.client.downloadReport(ctx, &r.amazonRequest, filePath)
	if err != nil {
		return err
	}
	return nil
}

/*
// Response is TSV data:

"_GET_FLAT_FILE_ORDERS_DATA_"
Index Legend:
[0] = order-id
[1] = order-item-id
[2] = purchase-date
[3] = payments-date
[4] = buyer-email
[5] = buyer-name
[6] = buyer-phone-number
[7] = sku
[8] = product-name
[9] = quantity-purchased
[10] = currency
[11] = item-price
[12] = item-tax
[13] = shipping-price
[14] = shipping-tax
[15] = ship-service-level
[16] = recipient-name
[17] = ship-address-1
[18] = ship-address-2
[19] = ship-address-3
[20] = ship-city
[21] = ship-state
[22] = ship-postal-code
[23] = ship-country
[24] = ship-phone-number
[25] = delivery-start-date
[26] = delivery-end-date
[27] = delivery-time-zone
[28] = delivery-Instructions
[29] = sales-channel
[30] = is-business-order
[31] = purchase-order-number
[32] = price-designationROW
*/
