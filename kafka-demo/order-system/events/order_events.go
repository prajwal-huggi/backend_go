package events

// =============================================================================
// SHARED EVENT DEFINITIONS
// =============================================================================
// In real microservice systems, you typically share a small "contract"
// package (or use a schema registry with Avro/Protobuf) so that the
// producer and ALL consumers agree on the exact shape of the message.
// Here we keep it simple with a Go struct + JSON, which is perfectly fine
// for learning and for many real systems too.
// =============================================================================

import (
	"encoding/json"
	"time"
)

// OrderCreated represents the event published whenever a new order is placed.
// This is the "contract" between the Order API (producer) and any service
// that cares about new orders (Inventory, Email, Analytics, etc.)
type OrderCreated struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	Items      []Item    `json:"items"`
	CreatedAt  time.Time `json:"created_at"`
}

type Item struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// Topic name used across the whole system - defining it once here avoids
// typos like "order-created" vs "orders-created" across services.
const OrderCreatedTopic = "order-created"

// ToJSON serializes the event for sending over Kafka.
func (o OrderCreated) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

// FromJSON deserializes a Kafka message value back into an OrderCreated event.
func FromJSON(data []byte) (OrderCreated, error) {
	var o OrderCreated
	err := json.Unmarshal(data, &o)
	return o, err
}
