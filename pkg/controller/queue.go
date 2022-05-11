package controller

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jjzcru/djinn/plugin"
	nats "github.com/nats-io/nats.go"
)

const MAX_RETRY_CONN = 5

type queue struct {
	nc      *nats.Conn
	plugins map[string]plugin.Plugin
}

func NewQueue(plugins map[string]plugin.Plugin) (queue, error) {
	q := queue{
		plugins: plugins,
	}
	nc, err := getQueueConn()

	if err != nil {
		return q, err
	}

	q.nc = nc
	return q, nil
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
