package models

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	producer *kafka.Producer
)

func init() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	producer = p
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		panic(err)
	}
}

func (baseModel *BaseModel) PublishEvents(message []byte, topic string) error {
	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return err
}
