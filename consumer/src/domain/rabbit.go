package domain

type RabbitMessage struct {
	Header      string `json:"header"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Status      string `json:"status"`
}

type RabbitMQRepository interface {
	ConsumeMessages(queueName string, handler func(message RabbitMessage) error) error
}