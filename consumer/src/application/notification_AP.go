package application

import (
	"fmt"
	"api/consumer/src/domain"
)

type NotificationAppService struct {
	notificationService domain.NotificationService
	tokenRepo           domain.TokenRepository
}

func NewNotificationAppService(notificationService domain.NotificationService, tokenRepo domain.TokenRepository) *NotificationAppService {
	return &NotificationAppService{
		notificationService: notificationService,
		tokenRepo:          tokenRepo,
	}
}

func (s *NotificationAppService) Subscribe(token string) (int64, error) {
	id, err := s.tokenRepo.SaveToken(token)
	if err != nil {
		return 0, fmt.Errorf("error guardando en BD: %v", err)
	}

	// Enviar notificación
	err = s.notificationService.SendNotification(token)
	if err != nil {
		return id, fmt.Errorf("error enviando notificación: %v", err)
	}

	return id, nil
}

func (s *NotificationAppService) SendRabbitNotification(header, description, image string) error {
	tokens, err := s.tokenRepo.GetAllTokens()
	if err != nil {
		return fmt.Errorf("error obteniendo tokens de usuarios: %v", err)
	}

	for _, token := range tokens {
		err := s.notificationService.SendFCMNotification(token, header, description, image)
		if err != nil {
			fmt.Printf("Error enviando notificación a %s: %v\n", token, err)
		}
	}

	return nil
}