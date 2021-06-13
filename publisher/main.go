package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"

	"github.com/streadway/amqp"
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

	// Create a new RabbitMQ Connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	panicIfNeeded(err)
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	panicIfNeeded(err)
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	panicIfNeeded(err)

	// Create a new Fiber instance.
	app := fiber.New()

	// Add middleware.
	app.Use(
		logger.New(), // add simple logger 
	)

	// Definer routes.
	routes(app, channelRabbitMQ)

	// Start Fiber API server.
	log.Fatal(app.Listen(":8080"))
}

func routes(app *fiber.App, mq *amqp.Channel) {
	// Add route for send message to Service 1.
	app.Post("/send", func(c *fiber.Ctx) error {
		// Create a message to publish.
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}

		// Attempt to publish a message to the queue.
		if err := mq.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		); err != nil {
			return err
		}

		return nil
	})
}
