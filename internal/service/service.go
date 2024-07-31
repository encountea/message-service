package service

import (
	"github.com/encountea/message-service/internal/kafka"
	"github.com/encountea/message-service/internal/models"
	"github.com/encountea/message-service/internal/repository"
)

type Service struct {
	repo     *repository.Repository
	producer *kafka.Producer
}

func NewService(repo *repository.Repository, producer *kafka.Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) ProcessMessage(msg models.Message) error {
	if err := s.repo.SaveMessage(msg); err != nil {
		return err
	}

	if err := s.producer.SendMessage(msg.Content); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetStats() (map[string]int, error) {
	count, err := s.repo.GetProcessedCount()
	if err != nil {
		return nil, err
	}

	return map[string]int{"processed_messages": count}, nil
}
