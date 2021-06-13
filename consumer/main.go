package main

import (
	"github.com/streadway/amqp"
	"log"
)

func panicIfNeeded(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Define RabbitMQ server URL.
	// amqpServerURL := "amqp://guest:guest@message-broker:5672/"
	amqpServerURL := "amqp://guest:guest@localhost:5672/"

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	panicIfNeeded(err)
	defer connectRabbitMQ.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	panicIfNeeded(err)
	defer channelRabbitMQ.Close()

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		"QueueService1", // queue name
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")
	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf("-> Received message: %s\n", message.Body)
		}
	}()

	<-forever
}