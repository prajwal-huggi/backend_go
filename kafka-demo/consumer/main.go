package main

// =============================================================================
// BASIC KAFKA CONSUMER (single reader, manually pinned to a partition)
// =============================================================================
// This is the SIMPLEST possible consumer - it reads directly from ONE
// partition. This is NOT how you'd normally do it in production (you'd
// use a Consumer Group, shown in the next file), but it's the best way
// to first understand the raw mechanics: offsets, partitions, polling.
// =============================================================================

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// -------------------------------------------------------------------
	// STEP 1: Create a "Reader" pinned to ONE specific partition.
	// -------------------------------------------------------------------
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"}, // can list multiple brokers
		Topic:     "order-created",
		Partition: 0, // <-- we are reading ONLY partition 0, manually

		// MinBytes/MaxBytes control batching behavior - Kafka will wait
		// until it has at least MinBytes of data (or a timeout), to avoid
		// making a network call per single tiny message.
		MinBytes: 1,    // 1 byte - don't wait for a big batch (good for learning/demo)
		MaxBytes: 10e6, // 10MB - max size of a single fetch
	})
	defer reader.Close()

	// -------------------------------------------------------------------
	// STEP 2: Tell it WHERE to start reading from.
	// -------------------------------------------------------------------
	// kafka.FirstOffset = start from the very beginning of the partition
	// (read all historical messages still within the retention period).
	// kafka.LastOffset = only read NEW messages from now on.
	reader.SetOffset(kafka.FirstOffset)

	log.Println("Starting to consume messages from partition 0...")

	// -------------------------------------------------------------------
	// STEP 3: Read messages in an infinite loop (this is normal! Kafka
	// consumers are typically long-running processes that just keep polling).
	// -------------------------------------------------------------------
	for {
		// ReadMessage blocks (waits) until a new message is available.
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		log.Printf(
			"📨 Received: partition=%d offset=%d key=%s value=%s",
			msg.Partition, msg.Offset, string(msg.Key), string(msg.Value),
		)

		// NOTE: With this simple Reader, the offset is committed
		// automatically as you read (because we used ReadMessage, which
		// auto-advances). In the Consumer Group example, we'll see how
		// offset committing is normally handled more deliberately.
	}
}
