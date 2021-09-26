package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn

func init(){

	var natsURL = nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}
	// Connect to the NATS server.
	NatsConn, _ = nats.Connect(natsURL, nats.Timeout(5*time.Second))
}


func main() {

	_,_= NatsConn.QueueSubscribe("message","queue1", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
	})

	fmt.Println("Press the Enter Key to terminate the console screen!")
	_, err := fmt.Scanln()
	if err != nil {
		return 
	} // wait for Enter Key
}
