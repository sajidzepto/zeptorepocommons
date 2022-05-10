package zeptorepocommons

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type BrokerProperties struct {
	topic   string
	address string
}

var (
	conn             *kafka.Conn
	brokerProperties *BrokerProperties
)

func InitializePublisher(properties *BrokerProperties) {
	if properties == nil {
		properties = &BrokerProperties{
			topic:   "default",
			address: "localhost:9092",
		}
	}
	con, err := kafka.DialLeader(context.Background(), "tcp", properties.address, properties.topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	conn = con
}

func (baseModel *BaseModel) PublishEvents(message []byte) error {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.WriteMessages(
		kafka.Message{Value: message},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	return err
}
