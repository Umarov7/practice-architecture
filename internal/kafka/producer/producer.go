package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type IKafkaProducer interface {
	Produce(ctx context.Context, topic string, msg []byte) error
	Close()
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) IKafkaProducer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{writer: w}
}

func (k *KafkaProducer) Produce(ctx context.Context, topic string, msg []byte) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: msg,
	})
}

func (k *KafkaProducer) Close() {
	k.writer.Close()
}
