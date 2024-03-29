package kafka

import (
	"context"
	"log"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafka.Reader
}

func NewKafkaReader() *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "topic_test",
		GroupID: "group",
	})

	return &Reader{
		Reader: reader,
	}
}

func (k *Reader) FetchMessage(ctx context.Context, messages chan<- kafka.Message) error {
	for {
		message, err := k.Reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case messages <- message:
			log.Printf("message fetched and sent to a channel: %v \n", string(message.Value))
		}
	}
}

func (k *Reader) CommitMessages(ctx context.Context, messageCommitChan <-chan kafka.Message) error {
	for {
		select {
		case <-ctx.Done():
		case msg := <-messageCommitChan:
			err := k.Reader.CommitMessages(ctx, msg)
			if err != nil {
				return errors.Wrap(err, "Reader.CommitMessages")
			}
			log.Printf("committed an msg: %v \n", string(msg.Value))
		}
	}
}