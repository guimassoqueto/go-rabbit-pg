package rabbitmq

import (
	"encoding/json"
	"grp/helpers"
	"grp/scraper"

	//"grp/scraper"
	"grp/types"
	"grp/variables"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

func Receive() {
	conn, err := ampq.Dial(variables.RABBITMQ_URL)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		variables.RABBITMQ_RECEIVE_QUEUE, // name
		false,                            // durable
		false,                            // delete when unusued
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  //exclusive
		false,  // no-local
		false,  //no-wait
		nil,    // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println("New item(s) received")

			var message types.RabbitMQMessage
			err := json.Unmarshal(d.Body, &message)

			if err != nil {
				d.Ack(false)
				log.Panicf("%s: %s", "Failed to parse the rabbitmq message.", err)
			}
			scraper.Scrap(message.AmazonColly)

			d.Ack(false)
		}
	}()

	log.Printf("Waiting for new messages [queue: amazon-colly]. To exit press CTRL+C")

	<-forever
}
