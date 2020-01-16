package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"net/url"
)

// ListOrderItemsRequest lists the line items associated with an order
type ListOrderItemsRequest struct {
	amazonRequest
}

// Do sends request to amazonMWS reports API
func (r *ListOrderItemsRequest) Do(ctx context.Context) (*ListOrderItemsResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &ListOrderItemsResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoNext sends gets next page of amazon reports if NextToken != ""
func (r *ListOrderItemsRequest) DoNext(ctx context.Context, nextToken string) (*ListOrderItemsByNextTokenResponse, error) {
	nextReq := &amazonRequest{
		client:   r.client,
		endpoint: r.endpoint,
		method:   r.method,
		params: url.Values{
			"Action":    {"ListOrderItemsByNextToken"},
			"Version":   r.params["Version"],
			"NextToken": {nextToken},
		},
	}

	respBytes, err := r.client.callAPI(ctx, nextReq)
	if err != nil {
		return nil, err
	}

	xmlResponse := &ListOrderItemsByNextTokenResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// DoAll calls multiple page response if necessary
func (r *ListOrderItemsRequest) DoAll(ctx context.Context) ([]OrderItem, error) {
	lineItems := []OrderItem{}
	firstResp, err := r.Do(ctx)
	if err != nil {
		return nil, err
	}
	lineItems = append(lineItems, firstResp.ListOrderItemsResult.OrderItems.OrderItem...)

	nextToken := firstResp.ListOrderItemsResult.NextToken
	hasNext := (nextToken != "")
	for hasNext {
		nextResp, err := r.DoNext(ctx, nextToken)
		if err != nil {
			return lineItems, err
		}
		lineItems = append(lineItems, nextResp.ListOrderItemsByNextTokenResult.OrderItems.OrderItem...)

		nextToken = nextResp.ListOrderItemsByNextTokenResult.NextToken
		hasNext = (nextToken != "")
	}

	return lineItems, nil
}

// ListOrderItemsResponse contains line item data
type ListOrderItemsResponse struct {
	XMLName              xml.Name `xml:"ListOrderItemsResponse"`
	Xmlns                string   `xml:"xmlns,attr"`
	ListOrderItemsResult struct {
		NextToken     string `xml:"NextToken"`
		AmazonOrderID string `xml:"AmazonOrderId"`
		OrderItems    struct {
			OrderItem []OrderItem `xml:"OrderItem"`
		} `xml:"OrderItems"`
	} `xml:"ListOrderItemsResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// ListOrderItemsByNextTokenResponse contains next page of line item data
type ListOrderItemsByNextTokenResponse struct {
	XMLName                         xml.Name `xml:"ListOrderItemsByNextTokenResponse"`
	Xmlns                           string   `xml:"xmlns,attr"`
	ListOrderItemsByNextTokenResult struct {
		NextToken     string `xml:"NextToken"`
		AmazonOrderID string `xml:"AmazonOrderId"`
		OrderItems    struct {
			OrderItem []OrderItem `xml:"OrderItem"`
		} `xml:"OrderItems"`
	} `xml:"ListOrderItemsByNextTokenResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// OrderItem contains order data - uncomment as needed
type OrderItem struct {
	SellerSKU       string `xml:"SellerSKU"`
	QuantityOrdered int    `xml:"QuantityOrdered"`
	OrderItemID     string `xml:"OrderItemId"`
	// ASIN                string `xml:"ASIN"`
	// BuyerCustomizedInfo struct {
	// 	CustomizedURL string `xml:"CustomizedURL"`
	// } `xml:"BuyerCustomizedInfo"`
	// Title           string `xml:"Title"`
	// QuantityShipped string `xml:"QuantityShipped"`
	// ProductInfo     struct {
	// 	NumberOfItems string `xml:"NumberOfItems"`
	// } `xml:"ProductInfo"`
	// PointsGranted struct {
	// 	PointsNumber        string `xml:"PointsNumber"`
	// 	PointsMonetaryValue struct {
	// 		CurrencyCode string `xml:"CurrencyCode"`
	// 		Amount       string `xml:"Amount"`
	// 	} `xml:"PointsMonetaryValue"`
	// } `xml:"PointsGranted"`
	// ItemPrice struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"ItemPrice"`
	// ShippingPrice struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"ShippingPrice"`
	// ScheduledDeliveryEndDate   string `xml:"ScheduledDeliveryEndDate"`
	// ScheduledDeliveryStartDate string `xml:"ScheduledDeliveryStartDate"`
	// CODFee                     struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"CODFee"`
	// CODFeeDiscount struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"CODFeeDiscount"`
	// IsGift          string `xml:"IsGift"`
	// IsTransparency  string `xml:"IsTransparency"`
	// GiftMessageText string `xml:"GiftMessageText"`
	// GiftWrapPrice   struct {
	// 	CurrencyCode string `xml:"CurrencyCode"`
	// 	Amount       string `xml:"Amount"`
	// } `xml:"GiftWrapPrice"`
	// GiftWrapLevel    string `xml:"GiftWrapLevel"`
	// PriceDesignation string `xml:"PriceDesignation"`
	// PromotionIds     struct {
	// 	PromotionId string `xml:"PromotionId"`
	// } `xml:"PromotionIds"`
	// ConditionId        string `xml:"ConditionId"`
	// ConditionSubtypeId string `xml:"ConditionSubtypeId"`
	// ConditionNote      string `xml:"ConditionNote"`
	// TaxCollection      struct {
	// 	Model            string `xml:"Model"`
	// 	ResponsibleParty string `xml:"ResponsibleParty"`
	// } `xml:"TaxCollection"`
}
