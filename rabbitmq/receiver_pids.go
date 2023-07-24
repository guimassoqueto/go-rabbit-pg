package rabbitmq

import (
	"grp/helpers"
	"grp/scraper"
	"grp/variables"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func ReceiverPids() {
	conn, err := ampq.Dial(variables.RABBITMQ_URL)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		variables.RABBITMQ_MAIN_QUEUE, // name
		false, // durable
		false, // delete when unusued
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	
	msgs, err := ch.Consume(
		q.Name, // queue
		"", // consumer
		true, // auto-ack
		false, //exclusive
		false, // no-local
		false, //no-wait
		nil, // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			message := string(d.Body[:])
			log.Printf("Received a message: %s", message)
			pidsArray := helpers.StringifiedArrayToArray(message)
			scraper.GoColly(pidsArray)
		}
	}()

	log.Printf(" [*] Waiting for messages. to exit press CTRL+C")

	<-forever
}