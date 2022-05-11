package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	nats "github.com/nats-io/nats.go"
)

type server struct {
  nc *nats.Conn
}

func (s server) baseRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Basic NATS based microservice example v0.0.1")
}

func (s server) healthz(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "OK")
}

func main() {
  var s server
  var err error
  uri := os.Getenv("NATS_URI")

	nc, err := nats.Connect(uri)

  _, err = nc.Subscribe("tasks", func(msg *nats.Msg){
      fmt.Printf("msg: %v\n", string(msg.Data))
  })

  if err != nil {
	 log.Fatal("Error establishing connection to NATS:", err)
  }

  fmt.Println("Connected to NATS at:", s.nc.ConnectedUrl())
  http.HandleFunc("/", s.baseRoot)
  http.HandleFunc("/healthz", s.healthz)

  fmt.Println("Server listening on port 8081...")
  if err := http.ListenAndServe(":8081", nil); err != nil {
	log.Fatal(err)
  }
}