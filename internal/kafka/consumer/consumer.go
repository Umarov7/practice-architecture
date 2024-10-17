package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type IKafkaConsumer interface {
	Consume(handler func(message []byte)) error
	Close()
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker []string, topic string) IKafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: broker,
				Topic:   topic,
			},
		),
	}
}

func (k *KafkaConsumer) Consume(handler func(message []byte)) error {
	for {
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		handler(m.Value)
	}
}

func (k *KafkaConsumer) Close() {
	k.reader.Close()
}
