package notification

import (
	"encoding/json"
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

type NotificationMessage struct {
	NotificationMethod string `json:"notificationMethod"`
	Rcpt               string `json:"rcpt"`
	Message            string `json:"message"`
}

var (
	client *zmq.Socket
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
func Connect() {
	fmt.Println("Connecting")
	client, _ = zmq.NewSocket(zmq.PUSH)
	client.Connect("tcp://127.0.0.1:3010")
}

func Send(not NotificationMessage) {
	b, err := json.Marshal(not)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.Send(string(b), 0)
}
