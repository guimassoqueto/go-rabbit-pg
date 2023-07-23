package main

import (
	"fmt"
	"grp/rabbitmq"
)

func main() {
	var i = 0
	for i < 10000 {
		rabbitmq.Send(fmt.Sprintf("Message #%d", i))
		i += 1
	}
	rabbitmq.Receive()
}