package main

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	service := &RabbitMQService{
		connection: conn,
		channel:    ch,
	}

	// Declarar exchange y queue
	err = service.setupQueues()
	if err != nil {
		service.Close()
		return nil, err
	}

	log.Println("✅ Connected to RabbitMQ")
	return service, nil
}

func (r *RabbitMQService) setupQueues() error {
	// Declarar exchange
	err := r.channel.ExchangeDeclare(
		"hotel.events", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// Declarar queue para el servicio de búsqueda
	_, err = r.channel.QueueDeclare(
		"hotel.search.updates", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue al exchange
	err = r.channel.QueueBind(
		"hotel.search.updates", // queue name
		"hotel.*",              // routing key
		"hotel.events",         // exchange
		false,
		nil,
	)

	return err
}

func (r *RabbitMQService) PublishMessage(routingKey string, message []byte) error {
	return r.channel.Publish(
		"hotel.events", // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQService) ConsumeMessages(queueName string, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true) // Requeue message
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQService) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}