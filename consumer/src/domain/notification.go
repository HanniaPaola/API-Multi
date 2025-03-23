package domain

type Notification struct {
	Token       string
	Header      string
	Description string
	Image       string
}

type NotificationService interface {
	SendNotification(token string) error
	SendFCMNotification(token, header, description, image string) error
}

type TokenRepository interface {
	SaveToken(token string) (int64, error)
	GetAllTokens() ([]string, error)
}