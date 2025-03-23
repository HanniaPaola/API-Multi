package infraestructure

import (
	"api/consumer/src/domain"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQRepositoryImpl struct {
	conn *amqp.Connection
}

func NewRabbitMQRepositoryImpl(url string) (*RabbitMQRepositoryImpl, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}
	return &RabbitMQRepositoryImpl{conn: conn}, nil
}

func (r *RabbitMQRepositoryImpl) ConsumeMessages(queueName string, handler func(message domain.RabbitMessage) error) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("error abriendo un canal: %v", err)
	}
	defer ch.Close()

	// Asegurar que la cola existe
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando la cola: %v", err)
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("error al consumir mensajes: %v", err)
	}

	// Procesar mensajes
	go func() {
		for msg := range msgs {
			var rabbitMsg domain.RabbitMessage
			if err := json.Unmarshal(msg.Body, &rabbitMsg); err != nil {
				log.Printf("Error decodificando mensaje: %v", err)
				continue
			}

			if err := handler(rabbitMsg); err != nil {
				log.Printf("Error procesando mensaje: %v", err)
			}
		}
	}()

	log.Println("Escuchando mensajes de RabbitMQ...")
	return nil
}