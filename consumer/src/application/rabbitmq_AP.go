package application

import (
	"fmt"
	"api/consumer/src/domain"
)

type RabbitMQAppService struct {
	rabbitRepo domain.RabbitMQRepository
	notificationService domain.NotificationService
}

func NewRabbitMQAppService(rabbitRepo domain.RabbitMQRepository, notificationService domain.NotificationService) *RabbitMQAppService {
	return &RabbitMQAppService{
		rabbitRepo:         rabbitRepo,
		notificationService: notificationService,
	}
}

func (s *RabbitMQAppService) StartConsuming(queueName string) error {
	err := s.rabbitRepo.ConsumeMessages(queueName, func(message domain.RabbitMessage) error {
		fmt.Println("Mensaje recibido de RabbitMQ:", message)

		err := s.notificationService.SendFCMNotification("token_del_usuario", message.Header, message.Description, message.Image)
		if err != nil {
			return fmt.Errorf("error enviando notificación: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error consumiendo mensajes: %v", err)
	}

	return nil
}