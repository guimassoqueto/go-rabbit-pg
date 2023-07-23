package rabbitmq

import (
	"context"
	"log"
	"time"

	"grp/helpers"
	"grp/variables"

	ampq "github.com/rabbitmq/amqp091-go"
)


func Send(message string) {
	conn, err := ampq.Dial(variables.RABBITMQ_URL)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()


	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		variables.RABBITMQ_MAIN_QUEUE, // name
		false, // durable
		false, // delete when unusues
		false, // exclusive
		false, // no-wait
		nil,	 // argumenst
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"", // exchange
		q.Name, // routing key
		false, // mandatory
		false, // inmediate
		ampq.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
	)
	helpers.FailOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s\n", message)
}