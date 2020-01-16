package amazonmwsapi

import (
	"context"
	"encoding/xml"
)

// GetOrderRequest holds request data for GetOrder call
type GetOrderRequest struct {
	amazonRequest
}

// Do sends request to Amazon API
func (r *GetOrderRequest) Do(ctx context.Context) (*GetOrderResponse, error) {
	respBytes, err := r.client.callAPI(ctx, &r.amazonRequest)
	if err != nil {
		return nil, err
	}

	xmlResponse := &GetOrderResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, r.client.parseAPIerrors(respBytes)
	}

	return xmlResponse, nil
}

// GetOrderResponse holds reponse data for ListOrders call
type GetOrderResponse struct {
	XMLName        xml.Name `xml:"GetOrderResponse"`
	Xmlns          string   `xml:"xmlns,attr"`
	GetOrderResult struct {
		Orders struct {
			Order []struct {
				AmazonOrderID          string `xml:"AmazonOrderId"`
				PurchaseDate           string `xml:"PurchaseDate"`
				LastUpdateDate         string `xml:"LastUpdateDate"`
				OrderStatus            string `xml:"OrderStatus"`
				FulfillmentChannel     string `xml:"FulfillmentChannel"`
				NumberOfItemsShipped   string `xml:"NumberOfItemsShipped"`
				NumberOfItemsUnshipped string `xml:"NumberOfItemsUnshipped"`
				PaymentMethod          string `xml:"PaymentMethod"`
				PaymentMethodDetails   struct {
					PaymentMethodDetail []string `xml:"PaymentMethodDetail"`
				} `xml:"PaymentMethodDetails"`
				MarketplaceID                string `xml:"MarketplaceId"`
				ShipmentServiceLevelCategory string `xml:"ShipmentServiceLevelCategory"`
				OrderType                    string `xml:"OrderType"`
				EarliestShipDate             string `xml:"EarliestShipDate"`
				LatestShipDate               string `xml:"LatestShipDate"`
				IsBusinessOrder              string `xml:"IsBusinessOrder"`
				IsPrime                      string `xml:"IsPrime"`
				IsPremiumOrder               string `xml:"IsPremiumOrder"`
			} `xml:"Order"`
		} `xml:"Orders"`
	} `xml:"GetOrderResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}
