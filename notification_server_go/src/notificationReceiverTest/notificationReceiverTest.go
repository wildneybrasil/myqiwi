package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Println("Connecting")
	server, _ := zmq.NewSocket(zmq.PUSH)
	server.Connect("tcp://127.0.0.1:3010")


   server.Send(`{ "notificationMethod":"push_ios", "rcpt": "410b2eb81afcb693b760dc1381a001876450bb8149d91eb35829851fa415bf18", "message": "teste", "params":[ {"teste": "teste2"} ]}`,0);
//   server.Send("{ \"notificationMethod\":\"email\", \"rcpt\": \"fyy@mac.com\", \"message\": \"teste\", \"result\":\"queue\"}",0);
//   server.Send("{ \"notificationMethod\":\"sms\", \"rcpt\": \"+5511989288082\", \"message\": \"teste\"}",0);
	time.Sleep(1*time.Minute)

}
