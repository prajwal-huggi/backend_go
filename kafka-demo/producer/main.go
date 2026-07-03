package main

// =============================================================================
// BASIC KAFKA PRODUCER
// =============================================================================
// This program connects to Kafka and sends ("produces") a few messages
// to a topic called "order-created".
//
// Library used: github.com/segmentio/kafka-go
// Install it with: go get github.com/segmentio/kafka-go
// =============================================================================

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// -------------------------------------------------------------------
	// STEP 1: Create a "Writer". Think of this as your connection/client
	// that knows how to talk to the Kafka cluster and send messages.
	// -------------------------------------------------------------------
	writer := &kafka.Writer{
		// Addr: the address of ANY broker in the cluster. The client
		// will automatically discover all other brokers from this one.
		// This is called the "bootstrap server".
		Addr: kafka.TCP("localhost:9092"),

		// Topic: the name of the topic we want to send messages to.
		// If it doesn't exist yet, Kafka will auto-create it
		// (because we enabled KAFKA_AUTO_CREATE_TOPICS_ENABLE in docker-compose).
		Topic: "order-created",

		// Balancer: this decides WHICH PARTITION a message goes to,
		// when no key is provided (or sometimes based on the key).
		// kafka.LeastBytes picks the partition with the le~ast data
		// buffered - good for even load distribution.
		// (We'll use kafka.Hash{} a bit later when keys matter.)
		Balancer: &kafka.LeastBytes{},

		// RequiredAcks: how many brokers must confirm they've received
		// the message before we consider the "send" successful.
		// kafka.RequireAll = wait for ALL in-sync replicas to acknowledge.
		// This is the "durability vs speed" tradeoff knob.
		RequiredAcks: kafka.RequireAll,
	}

	// IMPORTANT: always close the writer when you're done, this flushes
	// any buffered messages and closes the network connection cleanly.
	defer writer.Close()

	// -------------------------------------------------------------------
	// STEP 2: Send some messages
	// -------------------------------------------------------------------
	for i := 1; i <= 5; i++ {
		orderID := fmt.Sprintf("order-%d", i)

		message := kafka.Message{
			// Key: used by Kafka to decide the partition (same key
			// always goes to the same partition -> keeps order).
			// Here we use orderID as the key.
			Key: []byte(orderID),

			// Value: the actual message content/payload.
			// In real systems, this is usually JSON.
			Value: []byte(fmt.Sprintf(
				`{"order_id": "%s", "amount": %d, "status": "created"}`,
				orderID, i*100,
			)),

			// Time: optional, defaults to now if not set.
			Time: time.Now(),
		}

		// context.Background() = no special cancellation/timeout rules.
		// In production you'd often use context.WithTimeout(...).
		err := writer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Fatalf("failed to write message: %v", err)
		}

		log.Printf("✅ Produced message: %s", message.Value)

		time.Sleep(500 * time.Millisecond) // just to slow it down so it's readable
	}

	log.Println("All messages sent successfully!")
}
