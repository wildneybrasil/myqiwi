// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package main

import (
	"email"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"log"
	"os"
	"notificationReceiver"
	"push_android"
	"push_ios"
	"errors"
	"sms"
	"strconv"
)

var (
	uri              = flag.String("uri", "amqp://guest:guest@" + os.Getenv("SERVER_RABBITMQ") + ":5672/", "AMQP URI")
	exchange         = flag.String("exchange", "iris.nofification-exchange", "Durable, non-auto-deleted AMQP exchange name")
	waitExchange     = flag.String("wait-exchange", "iris.nofification-wait-exchange", "Durable, non-auto-deleted AMQP exchange name")
	responseExchange = flag.String("response-exchange", "iris.nofification-response-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType     = flag.String("exchange-type", "topic", "Exchange type - direct|fanout|topic|x-custom")
	queue            = flag.String("queue", "iris.notification", "Ephemeral AMQP queue name")
	waitQueue        = flag.String("queue-wait", "iris.notification-wait", "Ephemeral AMQP queue name")
	responseQueue    = flag.String("queue-response", "iris.notification-response", "Ephemeral AMQP queue name")
	consumerTag      = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

func init() {
	flag.Parse()
}

func main() {

	err := NewConsumer(*uri, *exchange, *waitExchange, *responseExchange, *exchangeType, *queue, *waitQueue, *responseQueue, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("running forever")
	select {}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, waitExchange, responseExchange, exchangeType, queueName, waitQueueName, responseQueueName, ctag string) error {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("Worker dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		waitExchange, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		amqp.Table{},
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		responseExchange, // name of the exchange
		exchangeType,     // type
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	waitQueue, err := c.channel.QueueDeclare(
		waitQueueName, // name of the queue
		true,          // durable
		false,         // delete when usused
		false,         // exclusive
		false,         // noWait
		amqp.Table{"x-dead-letter-exchange": exchange}, // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	responseQueue, err := c.channel.QueueDeclare(
		responseQueueName, // name of the queue
		true,              // durable
		false,             // delete when usused
		false,             // exclusive
		false,             // noWait
		amqp.Table{}, // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		responseQueue.Name, // name of the queue
		"iris.notification.*",                // bindingKey
		responseExchange,   // sourceExchange
		false,              // noWait
		nil,                // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = c.channel.QueueBind(
		waitQueue.Name, // name of the queue
		"iris.notification.*",            // bindingKey
		waitExchange,   // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		"iris.notification.*",        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	go receiver.Work(exchange, *c.channel)

	for {
		deliveries, err := c.channel.Consume(
			queue.Name, // name
			c.tag,      // consumerTag,
			false,      // noAck
			false,      // exclusive
			false,      // noLocal
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			return fmt.Errorf("Queue Consume: %s", err)
		}
		dispatcher := NewDispatcher(deliveries, 3)
		dispatcher.Dispatch(func(delivery amqp.Delivery, id int) {

			js, err := simplejson.NewJson(delivery.Body)
			errCrititcal := false;
			
			if err == nil {
				notificationType := js.Get("notificationMethod").MustString()
				fmt.Println("TYPE: " + notificationType)
				switch notificationType {
				case "sms":
					if( len(js.Get("rcpt").MustString() ) <4){
						err = errors.New("SMS RCPT too small")
						errCrititcal = true;
					}
					if( err == nil ){
						err = sms.Send(js.Get("rcpt").MustString(), js.Get("message").MustString())
					}
					break
				case "push":
					fmt.Println("PUSH")
					if( len(js.Get("rcpt").MustString()) < 10 ){
						err = errors.New("PUSH RCPT too small")
						errCrititcal = true;
					}
					if( err == nil ){
						if( len(js.Get("rcpt").MustString()) == 64 ){
							err = ios.Send(js.Get("rcpt").MustString(), 0, js.Get("message").MustString(), js.Get("params").MustMap())
						} else {
							rcpt :=  js.Get("rcpt").MustString()[8:]
							fmt.Println(js.Get("rcpt").MustString());
							fmt.Println("TO" + rcpt )
							err = android.Send(rcpt, "", js.Get("message").MustString(), js.Get("params").MustMap())						
						}
					}
					break
				case "push_ios":
					err = ios.Send(js.Get("rcpt").MustString(), 0, js.Get("message").MustString(), js.Get("params").MustMap() )
					break
					
				case "push_android":
					title:=""
					if(js.Get("title") !=nil ){
						title = js.Get("title").MustString()
					}
					
					err = android.Send(js.Get("rcpt").MustString(), title, js.Get("message").MustString(), js.Get("params").MustMap())
					break
				case "email":
					if( len(js.Get("rcpt").MustString()) < 5 ){
						err = errors.New("EMAIL RCPT too small")
						errCrititcal = true;
					}
					if( err == nil ){
						err = email.Send(js.Get("rcpt").MustString(), js.Get("message").MustString())
						fmt.Println("Sending email", err)
					}
					break
				default:
					fmt.Println("Uknown type: '"+ notificationType );
				break;
				}

				if err == nil {
					// se tudo deu certo, coloque mensagem na fila de resposta
					if js.Get("result").MustString() == "store" {
						fmt.Println("******************")
						
						if err = c.channel.Publish(
							responseExchange, // publish to an exchange
							delivery.RoutingKey,      // routing to 0 or more queues
							false,            // mandatory
							false,            // immediate
							amqp.Publishing{
								Headers:         amqp.Table{"x-notificationStatus": "success"},
								ContentType:     "text/plain",
								ContentEncoding: "",
								Body:            []byte(delivery.Body),
								DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
								Priority:        0,               // 0-9
							},
						); err != nil {
							fmt.Println("Exchange Publish: %s", err)
						}
					}
				} else {
					if( errCrititcal ){
						if js.Get("result").MustString() == "store" {
							if err = c.channel.Publish(
								responseExchange, // publish to an exchange
								delivery.RoutingKey,      // routing to 0 or more queues
								false,            // mandatory
								false,            // immediate
								amqp.Publishing{
									Headers:         amqp.Table{"x-notificationStatus": "failed"},
									ContentType:     "text/plain",
									ContentEncoding: "",
									Body:            []byte(delivery.Body),
									DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
									Priority:        0,               // 0-9
								},
							); err != nil {
								fmt.Println("Exchange Publish: %s", err)
							}
						}						
					} else {
						interval := js.Get("retry_interval").MustString()
						max := js.Get("retry_max_count").MustInt()
	
						if interval == "" {
							interval = "60000"
						}
						if max == 0 {
							max = 1
						}
						count := 0
	
						if delivery.Headers["x-count"] != nil {
							count, _ = strconv.Atoi(delivery.Headers["x-count"].(string))
							count = count - 1
						} else {
							count = max
						}
	
						if count != -1 {
							if err = c.channel.Publish(
								waitExchange, // publish to an exchange
								delivery.RoutingKey,  // routing to 0 or more queues
								false,        // mandatory
								false,        // immediate
								amqp.Publishing{
									Headers:         amqp.Table{"x-count": strconv.Itoa(count)},
									ContentType:     "text/plain",
									ContentEncoding: "",
									Expiration:      interval,
									Body:            []byte(delivery.Body),
									DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
									Priority:        0,               // 0-9
								},
							); err != nil {
								fmt.Println("Exchange Publish: %s", err)
							}
						} else {
							// se tudo deu errado, coloque mensagem na fila de resposta
							if js.Get("result").MustString() == "store" {
								if err = c.channel.Publish(
									responseExchange, // publish to an exchange
									delivery.RoutingKey,      // routing to 0 or more queues
									false,            // mandatory
									false,            // immediate
									amqp.Publishing{
										Headers:         amqp.Table{"x-notificationStatus": "failed"},
										ContentType:     "text/plain",
										ContentEncoding: "",
										Body:            []byte(delivery.Body),
										DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
										Priority:        0,               // 0-9
									},
								); err != nil {
									fmt.Println("Exchange Publish: %s", err)
								}
							}
	
						}
					}

				}
				delivery.Ack(false)
			}
		})
		dispatcher.Wait()
	}

}
