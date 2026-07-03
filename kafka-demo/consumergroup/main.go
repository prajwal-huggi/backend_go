package main

// =============================================================================
// CONSUMER GROUP (the pattern you'll actually use in real applications)
// =============================================================================
// Instead of manually picking a partition, we join a "Consumer Group".
// Kafka automatically assigns partitions to us, and if we run multiple
// copies of this same program (with the same GroupID), Kafka will
// automatically SPLIT the partitions between them. This is how you scale
// horizontally.
//
// Try this: run this program twice (in 2 terminals) with the same GroupID
// while the topic has multiple partitions, then watch in kafka-ui how
// each instance gets assigned different partitions.
// =============================================================================

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	// -------------------------------------------------------------------
	// STEP 1: Create a Reader, but this time with a GroupID set.
	// Setting GroupID is THE switch that turns this into a Consumer
	// Group member, instead of a manual partition reader.
	// -------------------------------------------------------------------
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order-created",

		// GroupID: all consumers using the SAME GroupID form one
		// Consumer Group. Kafka tracks each group's progress (offsets)
		// SEPARATELY per group. This is the magic that lets multiple
		// independent services (e.g. "email-service", "inventory-service")
		// read the same topic without interfering with each other.
		GroupID: "inventory-service-group",

		MinBytes: 1,
		MaxBytes: 10e6,

		// CommitInterval controls how often offsets are auto-committed
		// back to Kafka. 0 means commit synchronously after every
		// ReadMessage call (safest, but slightly slower). For learning,
		// keeping it explicit (0) is best so behavior is predictable.
		CommitInterval: 0,
	})
	defer reader.Close()

	// -------------------------------------------------------------------
	// STEP 2: Set up graceful shutdown.
	// In real services, you NEVER want to just kill -9 a consumer;
	// you want it to finish processing the current message and exit
	// cleanly. This listens for Ctrl+C (SIGINT) or SIGTERM (used by
	// Docker/Kubernetes when stopping a container).
	// -------------------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Println("Shutdown signal received, finishing up...")
		cancel()
	}()

	log.Println("Consumer group member started, waiting for messages...")

	// -------------------------------------------------------------------
	// STEP 3: The main consume loop
	// -------------------------------------------------------------------
	for {
		// FetchMessage gets the next message WITHOUT auto-committing
		// the offset. This gives us control: we only commit AFTER
		// we've successfully processed the message. This is important
		// for "at-least-once" delivery guarantees (explained below the code).
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Context canceled, shutting down cleanly.")
				break
			}
			log.Printf("error fetching message: %v", err)
			break
		}

		// ---------------------------------------------------------------
		// STEP 4: "Process" the message.
		// In a real app, this is where your business logic goes:
		// update inventory, save to DB, call another API, etc.
		// ---------------------------------------------------------------
		if err := processOrderEvent(msg); err != nil {
			// In production, you'd typically:
			// - retry a few times
			// - then send to a "dead letter queue" (a separate topic
			//   for messages that failed processing, for later inspection)
			log.Printf("⚠️ failed to process message, skipping commit: %v", err)
			continue // don't commit - we'll re-read this message next time
		}

		// ---------------------------------------------------------------
		// STEP 5: Commit the offset - tell Kafka "I've successfully
		// handled this message, move my position forward."
		// ---------------------------------------------------------------
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %v", err)
		}

		log.Printf(
			"✅ Processed & committed: partition=%d offset=%d value=%s",
			msg.Partition, msg.Offset, string(msg.Value),
		)
	}

	log.Println("Consumer exited cleanly.")
}

// processOrderEvent simuƒtes business logic - e.g. decrementing stock
// in a real inventory system.
func processOrderEvent(msg kafka.Message) error {
	log.Printf("Processing order event: %s", string(msg.Value))
	// Pretend we update inventory in a database here.
	// Return an error here to simulate a failure and see the retry behavior.
	return nil
}
