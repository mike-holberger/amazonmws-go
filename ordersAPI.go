package amazonmwsapi

import (
	"fmt"
	"net/url"
)

const ordersAPIversion = "2013-09-01"

// OrdersAPI hold calls for getting order data
type OrdersAPI struct {
	client   *AmazonClient
	endpoint string
}

// NewOrdersAPI creates and configures new OrdersAPI object
func NewOrdersAPI(cli *AmazonClient) *OrdersAPI {
	api := &OrdersAPI{client: cli}
	api.endpoint = api.client.Region.Endpoint + "Orders/" + ordersAPIversion
	return api
}

// ListOrders gets a selection of orders
func (api *OrdersAPI) ListOrders() *ListOrdersRequest {
	return &ListOrdersRequest{
		amazonRequest{
			client:   api.client,
			params:   url.Values{"Action": {"ListOrders"}, "Version": {ordersAPIversion}},
			endpoint: api.endpoint,
			method:   "GET",
		},
	}
}

// ListOrderItems downloads a single report
func (api *OrdersAPI) ListOrderItems(orderID string) *ListOrderItemsRequest {
	return &ListOrderItemsRequest{
		amazonRequest{
			client:   api.client,
			endpoint: api.endpoint,
			params: url.Values{"Action": {"ListOrderItems"}, "Version": {ordersAPIversion},
				"AmazonOrderId": {orderID}},
			method: "GET",
		},
	}
}

// GetOrder gets a selection of orders
func (api *OrdersAPI) GetOrder(orderIDs []string) *GetOrderRequest {
	req := &GetOrderRequest{
		amazonRequest{
			client: api.client,
			params: url.Values{"Action": {"GetOrder"}, "Version": {ordersAPIversion}}, endpoint: api.endpoint,
			method: "GET",
		},
	}

	for i, id := range orderIDs {
		key := fmt.Sprintf("AmazonOrderId.Id.%d", (i + 1))
		req.params.Add(key, id)
	}

	return req
}
