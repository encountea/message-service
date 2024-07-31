package kafka

import (
	"github.com/IBM/sarama"
	"github.com/encountea/message-service/config"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	producer, err := sarama.NewSyncProducer(cfg.Brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, topic: cfg.Topic}, nil
}

func (p *Producer) SendMessage(msg string) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	})
	return err
}
