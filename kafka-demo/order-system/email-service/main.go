package main

// =============================================================================
// EMAIL SERVICE - reacts to OrderCreated events by sending confirmation emails
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

		// Different GroupID from inventory-service! This is THE reason
		// both services can independently read 100% of the messages.
		// If we accidentally used the SAME GroupID as inventory-service,
		// Kafka would treat them as the same group, and the partitions
		// would be SPLIT between them - each one would only see HALF
		// the messages. That's a classic beginner bug, so remember this.
		GroupID: "email-service-group",
	})
	defer reader.Close()

	log.Println("📧 Email service started, listening for orders...")

	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		order, err := events.FromJSON(msg.Value)
		if err != nil {
			log.Printf("bad message, skipping: %v", err)
			reader.CommitMessages(context.Background(), msg)
			continue
		}

		// --- business logic: send confirmation email ---
		log.Printf("📧 Sending confirmation email to customer=%s for order=%s ($%.2f)",
			order.CustomerID, order.OrderID, order.Amount)

		if err := reader.CommitMessages(context.Background(), msg); err != nil {
			log.Printf("commit failed: %v", err)
		}
	}
}
