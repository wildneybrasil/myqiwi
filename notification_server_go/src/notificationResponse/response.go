package main
 
import (
    "flag"
_    "fmt"
    "io/ioutil"
    "crypto/cipher"  
    "time" 
    "runtime"
_    "encoding/hex"
    zmq "github.com/pebbe/zmq4"
 
)
 


func Server() {
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()

	receiver.Connect("tcp://localhost:5557")
 
}
