package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/teran/linker/models"
	"github.com/teran/linker/repositories/stats"
	protov1 "github.com/teran/linker/repositories/stats/kafka/proto/v1"
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
	pr, err := protov1.NewRequest(req)
	if err != nil {
		return errors.Wrap(err, "error creating DTO")
	}

	payload, err := proto.Marshal(pr)
	if err != nil {
		return errors.Wrap(err, "error marshaling DTO")
	}

	_, _, err = r.producer.SendMessage(&sarama.ProducerMessage{
		Topic: r.topic,
		Value: sarama.ByteEncoder(payload),
	})
	return errors.Wrap(err, "error producing message")
}
