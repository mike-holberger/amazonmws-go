package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

// ListOrdersRequest holds request data for ListOrders call
type ListOrdersRequest struct {
	amazonRequest
}

// CreatedAfter ...
func (r *ListOrdersRequest) CreatedAfter(t time.Time) *ListOrdersRequest {
	xmlTime := XMLTimestamp(t)
	r.params.Add("CreatedAfter", xmlTime)
	return r
}

// CreatedBefore ...
func (r *ListOrdersRequest) CreatedBefore(t time.Time) *ListOrdersRequest {
	xmlTime := XMLTimestamp(t)
	r.params.Add("CreatedBefore", xmlTime)
	return r
}

// OrderStatus ...
func (r *ListOrdersRequest) OrderStatus(status []string) *ListOrdersRequest {
	for i, stat := range status {
		key := fmt.Sprintf("OrderStatus.Status.%d", (i + 1))
		r.params.Add(key, stat)
	}
	return r
}

// Do sends request to Amazon API
func (r *ListOrdersRequest) Do(ctx context.Context) (*ListOrdersResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &ListOrdersResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoNext sends gets next page of amazon reports if HasNext == true
func (r *ListOrdersRequest) DoNext(ctx context.Context, nextToken string) (*ListOrdersByNextTokenResponse, error) {
	nextReq := &amazonRequest{
		client:   r.client,
		endpoint: r.endpoint,
		method:   r.method,
		params: url.Values{
			"Action":    {"ListOrdersByNextToken"},
			"Version":   r.params["Version"],
			"NextToken": {nextToken},
		},
	}

	respBytes, err := r.client.callAPI(ctx, nextReq)
	if err != nil {
		return nil, err
	}

	xmlResponse := &ListOrdersByNextTokenResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoAll calls multiple page response if necessary
func (r *ListOrdersRequest) DoAll(ctx context.Context) ([]Order, error) {
	orders := []Order{}
	firstResp, err := r.Do(ctx)
	if err != nil {
		return nil, err
	}
	orders = append(orders, firstResp.ListOrdersResult.Orders.Order...)

	nextToken := firstResp.ListOrdersResult.NextToken
	hasNext := (nextToken != "")
	for hasNext {
		nextResp, err := r.DoNext(ctx, nextToken)
		if err != nil {
			return orders, err
		}
		orders = append(orders, nextResp.ListOrdersByNextTokenResult.Orders.Order...)

		nextToken = nextResp.ListOrdersByNextTokenResult.NextToken
		hasNext = (nextToken != "")
	}

	return orders, nil
}

// ListOrdersResponse holds reponse data for ListOrders call
type ListOrdersResponse struct {
	XMLName          xml.Name `xml:"ListOrdersResponse"`
	Xmlns            string   `xml:"xmlns,attr"`
	ListOrdersResult struct {
		NextToken         string `xml:"NextToken"`
		LastUpdatedBefore string `xml:"LastUpdatedBefore"`
		Orders            struct {
			Order []Order `xml:"Order"`
		} `xml:"Orders"`
	} `xml:"ListOrdersResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// ListOrdersByNextTokenResponse holds reponse data for ListOrders call
type ListOrdersByNextTokenResponse struct {
	XMLName                     xml.Name `xml:"ListOrdersByNextTokenResponse"`
	Text                        string   `xml:",chardata"`
	Xmlns                       string   `xml:"xmlns,attr"`
	ListOrdersByNextTokenResult struct {
		Text              string `xml:",chardata"`
		NextToken         string `xml:"NextToken"`
		LastUpdatedBefore string `xml:"LastUpdatedBefore"`
		Orders            struct {
			Text  string  `xml:",chardata"`
			Order []Order `xml:"Order"`
		} `xml:"Orders"`
	} `xml:"ListOrdersByNextTokenResult"`
	ResponseMetadata struct {
		Text      string `xml:",chardata"`
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// Order contains data on a single order
type Order struct {
	AmazonOrderID   string `xml:"AmazonOrderId"`
	BuyerEmail      string `xml:"BuyerEmail"`
	ShippingAddress struct {
		Name          string `xml:"Name"`
		AddressLine1  string `xml:"AddressLine1"`
		AddressLine2  string `xml:"AddressLine2"`
		AddressLine3  string `xml:"AddressLine3"`
		City          string `xml:"City"`
		StateOrRegion string `xml:"StateOrRegion"`
		PostalCode    string `xml:"PostalCode"`
		CountryCode   string `xml:"CountryCode"`
		Phone         string `xml:"Phone"`
		AddressType   string `xml:"AddressType"`
	} `xml:"ShippingAddress"`
	// PurchaseDate       time.Time `xml:"PurchaseDate"`
	// LastUpdateDate     string    `xml:"LastUpdateDate"`
	// OrderStatus        string    `xml:"OrderStatus"`
	// FulfillmentChannel string    `xml:"FulfillmentChannel"`
	// SalesChannel       string    `xml:"SalesChannel"`
	// OrderTotal struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"OrderTotal"`
	// NumberOfItemsShipped   string `xml:"NumberOfItemsShipped"`
	// NumberOfItemsUnshipped string `xml:"NumberOfItemsUnshipped"`
	// PaymentMethod          string `xml:"PaymentMethod"`
	// PaymentMethodDetails   struct {
	// 	PaymentMethodDetail string `xml:"PaymentMethodDetail"`
	// } `xml:"PaymentMethodDetails"`
	// MarketplaceID string `xml:"MarketplaceId"`
	// BuyerName     string `xml:"BuyerName"`
	// BuyerTaxInfo  struct {
	// 	CompanyLegalName   string `xml:"CompanyLegalName"`
	// 	TaxingRegion       string `xml:"TaxingRegion"`
	// 	TaxClassifications struct {
	// 		TaxClassification struct {
	// 			Text  string `xml:",chardata"`
	// 			Name  string `xml:"Name"`
	// 			Value string `xml:"Value"`
	// 		} `xml:"TaxClassification"`
	// 	} `xml:"TaxClassifications"`
	// } `xml:"BuyerTaxInfo"`
	// OrderType              string `xml:"OrderType"`
	// EarliestShipDate       string `xml:"EarliestShipDate"`
	// LatestShipDate         string `xml:"LatestShipDate"`
	// IsBusinessOrder        string `xml:"IsBusinessOrder"`
	// PurchaseOrderNumber    string `xml:"PurchaseOrderNumber"`
	// IsPrime                string `xml:"IsPrime"`
	// IsPremiumOrder         string `xml:"IsPremiumOrder"`
	// BuyerCounty            string `xml:"BuyerCounty"`
	// ShipServiceLevel       string `xml:"ShipServiceLevel"`
	// PaymentExecutionDetail struct {
	// 	PaymentExecutionDetailItem []struct {
	// 		Payment struct {
	// 			Text         string `xml:",chardata"`
	// 			Amount       string `xml:"Amount"`
	// 			CurrencyCode string `xml:"CurrencyCode"`
	// 		} `xml:"Payment"`
	// 		PaymentMethod string `xml:"PaymentMethod"`
	// 	} `xml:"PaymentExecutionDetailItem"`
	// } `xml:"PaymentExecutionDetail"`
	// ShipmentServiceLevelCategory string `xml:"ShipmentServiceLevelCategory"`
	// PromiseResponseDueDate       string `xml:"PromiseResponseDueDate"`
	// IsEstimatedShipDateSet       string `xml:"IsEstimatedShipDateSet"`
}
