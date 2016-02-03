package receiver

import (
	"flag"
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"
	zmq "github.com/pebbe/zmq4"
	"github.com/streadway/amqp"
)

var (
	reliable = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func publish(exchange string, routingKey string, channel amqp.Channel, body string) error {

	log.Printf("declared Exchange, publishing " + body)
	if err := channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func Work(exchange string, channel amqp.Channel) {
	fmt.Println("Receiver ")

	server, _ := zmq.NewSocket(zmq.PULL)
	server.Bind("tcp://*:3010")

	for {
		request, _ := server.Recv(0)

		js, err := simplejson.NewJson([]byte(request))

		routingKey := ""

		if err == nil {
			notificationType := js.Get("notificationMethod").MustString()
			fmt.Println("TYPE: " + notificationType)
			switch notificationType {
			case "sms":
				routingKey = "iris.notification.sms"
				break
			case "push":
				routingKey = "iris.notification.push"
				break
			case "push_ios":
				routingKey = "iris.notification.push_ios"
				break
			case "push_android":
				routingKey = "iris.notification.push_android"
				break
			case "email":
				routingKey = "iris.notification.email"
				break
			}

			fmt.Println("Received")

			if err := publish(exchange, routingKey, channel, request); err != nil {
				log.Fatalf("%s", err)
			}
		} else {
			fmt.Println(err)
		}
	}
}
