package amazonmwsapi

import (
	"context"
	"encoding/xml"
	"errors"
)

// SubmitFeedRequest holds request data for SubmitFeed call
type SubmitFeedRequest struct {
	amazonRequest
	feed *feed
}
type feed struct {
	XMLName                   xml.Name    `xml:"AmazonEnvelope"`
	Xsi                       string      `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string      `xml:"noNamespaceSchemaLocation,attr"`
	Header                    *feedHeader `xml:"Header"`
	MessageType               string      `xml:"MessageType"`

	// ALLOWS ONLY ONE OF THESE CATEGORIES, CORRESPONDING TO MessageType
	OrderAckMessages         []*feedOrderAckMessage         // `xml:"Message,omitempty"`
	OrderFulfillmentMessages []*feedOrderFulfillmentMessage // `xml:"Message,omitempty"`
	InventoryMessages        []*feedInventoryMessage        // `xml:"Message,omitempty"`
	PriceMessages            []*feedPriceMessage            // `xml:"Message,omitempty"`
}
type feedHeader struct {
	DocumentVersion    string `xml:"DocumentVersion"`
	MerchantIdentifier string `xml:"MerchantIdentifier"`
}

func (r *SubmitFeedRequest) newFeed(messageType string) *feed {
	return &feed{
		Xsi:                       "http://www.w3.org/2001/XMLSchema-instance",
		NoNamespaceSchemaLocation: "amzn-envelope.xsd",
		Header: &feedHeader{
			DocumentVersion:    "1.01",
			MerchantIdentifier: r.client.credentials.Merchant,
		},
		MessageType: messageType,
	}
}

func (f feed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{Name: xml.Name{Local: "AmazonEnvelope"}, Attr: []xml.Attr{
		xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: f.Xsi},
		xml.Attr{Name: xml.Name{Local: "xsi:noNamespaceSchemaLocation"}, Value: f.NoNamespaceSchemaLocation},
	}}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	err = e.EncodeElement(f.Header, xml.StartElement{Name: xml.Name{Local: "Header"}})
	if err != nil {
		return err
	}
	err = e.EncodeElement(f.MessageType, xml.StartElement{Name: xml.Name{Local: "MessageType"}})
	if err != nil {
		return err
	}

	// Avoids repeated tag 'Message' marshalling conflict
	switch {
	case len(f.OrderAckMessages) != 0:
		err = e.EncodeElement(f.OrderAckMessages, xml.StartElement{Name: xml.Name{Local: "Message"}})
	case len(f.OrderFulfillmentMessages) != 0:
		err = e.EncodeElement(f.OrderFulfillmentMessages, xml.StartElement{Name: xml.Name{Local: "Message"}})
	case len(f.InventoryMessages) != 0:
		err = e.EncodeElement(f.InventoryMessages, xml.StartElement{Name: xml.Name{Local: "Message"}})
	case len(f.PriceMessages) != 0:
		err = e.EncodeElement(f.PriceMessages, xml.StartElement{Name: xml.Name{Local: "Message"}})
	default:
		return errors.New("MarshalXML error: no feed messages")
	}
	if err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

// Do encodes XML feed, calculates MD5 sum, and submits to amazon feedsAPI
func (r *SubmitFeedRequest) Do(ctx context.Context) (*SubmitFeedResponse, error) {
	respBytes, err := r.client.submitFeed(ctx, r)
	if err != nil {
		return nil, err
	}

	xmlResponse := &SubmitFeedResponse{}
	err = xml.Unmarshal(respBytes, xmlResponse)
	if err != nil {
		return nil, err
	}

	return xmlResponse, nil
}

// SubmitFeedResponse obj
type SubmitFeedResponse struct {
	XMLName          xml.Name `xml:"SubmitFeedResponse"`
	Xmlns            string   `xml:"xmlns,attr"`
	SubmitFeedResult struct {
		FeedSubmissionInfo struct {
			FeedSubmissionID     string `xml:"FeedSubmissionId"`
			FeedType             string `xml:"FeedType"`
			SubmittedDate        string `xml:"SubmittedDate"`
			FeedProcessingStatus string `xml:"FeedProcessingStatus"`
		} `xml:"FeedSubmissionInfo"`
	} `xml:"SubmitFeedResult"`
	ResponseMetadata struct {
		RequestID string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

// OrderAcknowledgements creates XML feed body containing order acknowledgement data
func (r *SubmitFeedRequest) OrderAcknowledgements(orderIDs []string) *SubmitFeedRequest {
	r.params["FeedType"] = []string{"_POST_ORDER_ACKNOWLEDGEMENT_DATA_"}
	// create feed messages
	messages := make([]*feedOrderAckMessage, len(orderIDs))
	for i, ordID := range orderIDs {
		messages[i] = &feedOrderAckMessage{
			MessageID: (i + 1),
			OrderAck: &orderAcknowledgement{
				AmazonOrderID: ordID,
				StatusCode:    "Success",
			},
		}
	}

	// create feed obj and attach messages
	r.feed = r.newFeed("OrderAcknowledgement")
	r.feed.OrderAckMessages = messages
	return r
}

// Order Acknowledgement Feed structs
type feedOrderAckMessage struct {
	MessageID int                   `xml:"MessageID"`
	OrderAck  *orderAcknowledgement `xml:"OrderAcknowledgement"`
}
type orderAcknowledgement struct {
	AmazonOrderID string `xml:"AmazonOrderID"`
	StatusCode    string `xml:"StatusCode"`
}

// OrderFulfillmentFeed sends data such as shipping tracking to complete orders
func (r *SubmitFeedRequest) OrderFulfillmentFeed(orderFulfillmentData []*OrderFulfillment) *SubmitFeedRequest {
	r.params["FeedType"] = []string{"_POST_ORDER_FULFILLMENT_DATA_"}

	messages := make([]*feedOrderFulfillmentMessage, len(orderFulfillmentData))
	for i, data := range orderFulfillmentData {
		messages[i] = &feedOrderFulfillmentMessage{
			MessageID:        (i + 1),
			OrderFulfillment: data,
		}
	}

	// // create feed obj and attach messages
	r.feed = r.newFeed("OrderFulfillment")
	r.feed.OrderFulfillmentMessages = messages
	return r
}

type feedOrderFulfillmentMessage struct {
	MessageID        int               `xml:"MessageID"`
	OrderFulfillment *OrderFulfillment `xml:"OrderFulfillment"`
}

// OrderFulfillment contains order shipping data
type OrderFulfillment struct {
	AmazonOrderID   string           `xml:"AmazonOrderID"`
	FulfillmentDate string           `xml:"FulfillmentDate"`
	FulfillmentData *FulfillmentData `xml:"FulfillmentData"`
	Item            []Item           `xml:"Item"`
}

// FulfillmentData type
type FulfillmentData struct {
	CarrierCode           string `xml:"CarrierCode"`
	ShipperTrackingNumber string `xml:"ShipperTrackingNumber"`
}

// Item type
type Item struct {
	AmazonOrderItemCode string `xml:"AmazonOrderItemCode"`
	Quantity            int    `xml:"Quantity"`
}

// ReviseInventoryFeed creates XML feed body containing inventory revision data
func (r *SubmitFeedRequest) ReviseInventoryFeed(revisedItems []*FeedInventory) *SubmitFeedRequest {
	r.params["FeedType"] = []string{"_POST_INVENTORY_AVAILABILITY_DATA_"}

	// create feed messages
	messages := make([]*feedInventoryMessage, len(revisedItems))
	for i, item := range revisedItems {
		messages[i] = &feedInventoryMessage{
			MessageID:     (i + 1),
			OperationType: "Update",
			Inventory:     item,
		}
	}

	// create feed obj and attach messages
	r.feed = r.newFeed("Inventory")
	r.feed.InventoryMessages = messages

	return r
}

// Inventory Feed structs
type feedInventoryMessage struct {
	MessageID     int            `xml:"MessageID"`
	OperationType string         `xml:"OperationType"`
	Inventory     *FeedInventory `xml:"Inventory"`
}

// FeedInventory holds product inventory data
type FeedInventory struct {
	SKU      string `xml:"SKU"`
	Quantity int    `xml:"Quantity"`
}

// RevisePriceFeed creates XML feed body containing price revision data
func (r *SubmitFeedRequest) RevisePriceFeed(revisedItems []*RevisedPriceItem) *SubmitFeedRequest {
	r.params["FeedType"] = []string{"_POST_PRODUCT_PRICING_DATA_"}

	// create feed messages
	messages := make([]*feedPriceMessage, len(revisedItems))
	for i, item := range revisedItems {
		messages[i] = &feedPriceMessage{
			MessageID: (i + 1),
			Price: &feedPrice{
				SKU: item.SKU,
				StandardPrice: &StandardPrice{
					Amount:   item.Price,
					Currency: "USD",
				},
			},
		}
	}

	// create feed obj and attach messages
	r.feed = r.newFeed("Price")
	r.feed.PriceMessages = messages

	return r
}

// RevisedPriceItem request struct for adding revised items to feed messages
type RevisedPriceItem struct {
	SKU   string
	Price float64
}

// Price feed structs
type feedPriceMessage struct {
	MessageID int        `xml:"MessageID"`
	Price     *feedPrice `xml:"Price"`
}
type feedPrice struct {
	SKU           string         `xml:"SKU"`
	StandardPrice *StandardPrice `xml:"StandardPrice"`
}

// StandardPrice struct
type StandardPrice struct {
	Amount   float64 `xml:",chardata"`
	Currency string  `xml:"currency,attr"`
}
