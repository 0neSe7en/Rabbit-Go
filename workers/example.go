package workers

import (
	"github.com/Rabbit-Go/Base"
	"github.com/streadway/amqp"
	"log"
)

var Example = Base.WorkerConf{
	QName:     "Example",
	RoutingKs: []string{"#"},
	AutoAck:   false,
	Handler:   exampleHandler,
	Desc:      "Just to show how to write a worker in golang :)",
}

func exampleHandler(msg amqp.Delivery) {
	log.Printf(" [x] %s", msg.Body)
	msg.Ack(true)
}
