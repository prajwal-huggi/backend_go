package main

import (
	"fmt"
	"time"

	"github.com/prajwal-huggi/backend_go/internal/consumer"
	"github.com/prajwal-huggi/backend_go/internal/producer"
)

type Server struct {
	producer *producer.KafkaProducer
	consumer *consumer.KafkaConsumer
	msgCH chan string
}

func NewServer() *Server{
	msgCH:= make(chan string, 64)
	c, err:= consumer.NewKafkaConsumer(msgCH)

	if err!= nil{
		panic(err) 
	}

	return &Server{
		producer: producer.NewKafkaProducer(""),
		consumer: c,
		msgCH:  msgCH,
	}
}

func (s *Server) produceMsg(){
	ticker:= time.NewTicker(time.Second)
	defer ticker.Stop()
	id:= 0
	for t:= range ticker.C{
		msg:= fmt.Sprintf("hello from kafka, msgId= %d, ts= %s", id, t.Format("15:20:20"))
		s.producer.Produce(msg)
		id++
	}
}

func main(){
	s:= NewServer()
	s.produceMsg()
}