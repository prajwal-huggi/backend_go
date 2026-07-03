package main

// =============================================================================
// ORDER API - publishes OrderCreated events
// =============================================================================
// This simulates an HTTP API that, when a customer places an order, saves
// it to a database (not shown - keep it simple) AND publishes an event to
// Kafka so other services can react.
// =============================================================================

import (
	"context"
	"log"
	"time"

	"kafka-demo/order-system/events" // our shared contract package

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: events.OrderCreatedTopic,

		// kafka.Hash{} hashes the MESSAGE KEY to pick a partition.
		// We use this (instead of LeastBytes) because we WANT all events
		// for the same CustomerID to consistently land on the same
		// partition - useful if you ever care about per-customer ordering.
		Balancer: &kafka.Hash{},

		RequiredAcks: kafka.RequireAll,
	}
	defer writer.Close()

	// Simulate 3 customers placing orders
	orders := []events.OrderCreated{
		{
			OrderID:    "order-1001",
			CustomerID: "cust-A",
			Amount:     59.99,
			Items:      []events.Item{{SKU: "SKU-RED-SHOE", Quantity: 1}},
			CreatedAt:  time.Now(),
		},
		{
			OrderID:    "order-1002",
			CustomerID: "cust-B",
			Amount:     129.50,
			Items:      []events.Item{{SKU: "SKU-BLUE-JACKET", Quantity: 1}, {SKU: "SKU-CAP", Quantity: 2}},
			CreatedAt:  time.Now(),
		},
		{
			OrderID:    "order-1003",
			CustomerID: "cust-A", // same customer as order-1001 - same partition!
			Amount:     15.00,
			Items:      []events.Item{{SKU: "SKU-SOCKS", Quantity: 3}},
			CreatedAt:  time.Now(),
		},
	}

	for _, order := range orders {
		payload, err := order.ToJSON()
		if err != nil {
			log.Fatalf("failed to marshal order: %v", err)
		}

		msg := kafka.Message{
			Key:   []byte(order.CustomerID), // <-- partition key choice
			Value: payload,
		}

		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Fatalf("failed to write message: %v", err)
		}

		log.Printf("📤 Published OrderCreated: order_id=%s customer_id=%s", order.OrderID, order.CustomerID)
		time.Sleep(300 * time.Millisecond)
	}

	log.Println("Order API: all orders published.")
}
