package main

import (
	"os"
	"errors"
	"fmt"
	"log"
	
	"time"
	nats "github.com/nats-io/nats.go"
)

const MAX_RETRY_CONN = 5

func main() {
	nc, err := getQueueConn()
	if err != nil {
		os.Exit(1)
	}
	nc.Publish("command", []byte(`{
		"id": "d073d568f92a",
		"plugin": "com.lifx.bulbs",
		"command": "SET_COLOR",
		"payload": "{\"color\": \"#FFFF00\"}"
	}`))
	nc.Close()
	fmt.Println("Send command")
}

func getQueueConn() (*nats.Conn, error) {
	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn
	for i := 0; i < MAX_RETRY_CONN; i++ {
		nc, err = nats.Connect(uri)

		if err == nil {
			break
		} else {
			log.Fatal(err)
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error establishing connection to NATS: %s", err.Error()))
	}

	return nc, nil
}
