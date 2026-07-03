package main

// =============================================================================
// INVENTORY SERVICE - reacts to OrderCreated events by reducing stock
// =============================================================================

import (
	"context"
	"log"

	"kafka-demo/order-system/events"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   events.OrderCreatedTopic,

		// Its OWN group id - completely separate from email-service's group.
		// This means Inventory reads every single message independently,
		// regardless of what Email service has or hasn't read yet.
		GroupID: "inventory-service-group",
	})
	defer reader.Close()

	log.Println("📦 Inventory service started, listening for orders...")

	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		order, err := events.FromJSON(msg.Value)
		if err != nil {
			log.Printf("bad message, skipping: %v", err)
			reader.CommitMessages(context.Background(), msg) // skip poison message
			continue
		}

		// --- business logic: reduce stock for each item ---
		for _, item := range order.Items {
			log.Printf("📦 Reducing stock: sku=%s qty=%d (order=%s)", item.SKU, item.Quantity, order.OrderID)
		}

		if err := reader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("commit failed: %v", err)
		}
	}
}
