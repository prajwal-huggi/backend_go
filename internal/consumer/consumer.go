package consumer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/prajwal-huggi/backend_go/internal/shared"
)

type KafkaConsumer struct{
	consumer *kafka.Consumer
	topic string
	msgCH chan<- string
} 

func NewKafkaConsumer(msgCH chan<- string)(*KafkaConsumer, error){
	cfg:= shared.NewKafkaConfig()
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
		"group.id":          cfg.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	// defer c.Close()

	err = c.SubscribeTopics([]string{cfg.Topic}, nil)

	if err != nil {
		panic(err)
	}
	consumer:= &KafkaConsumer{
		consumer: c,
		topic: cfg.Topic,
		msgCH: msgCH,
	}

	go consumer.readMsgLoop()

	return consumer, nil

}

func (c* KafkaConsumer) readMsgLoop() {
	for {
		msg, err := c.consumer.ReadMessage(time.Second)
		if err != nil && !err.(kafka.Error).IsTimeout(){
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			continue
		}
		
		payload:= msg.Value
		c.msgCH <- string(payload)
	}
}