package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/stats"
)

type repository struct {
	topic    string
	producer sarama.SyncProducer
}

func New(producer sarama.SyncProducer, topic string) stats.Repository {
	return &repository{
		topic:    topic,
		producer: producer,
	}
}

func (r *repository) LogRequest(ctx context.Context, req models.Request) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "error marshaling request")
	}

	_, _, err = r.producer.SendMessage(&sarama.ProducerMessage{
		Topic: r.topic,
		Value: sarama.ByteEncoder(payload),
	})
	return errors.Wrap(err, "error producing message")
}
