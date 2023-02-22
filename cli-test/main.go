package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"flag"

	"github.com/nats-io/nats.go"
)

func main() {
	natsurl := flag.String("natsurl", "127.0.0.1:4222", "nats url. Default is 127.0.0.1:4222")
	devicename := flag.String("devicename", "ceos1", "Name of the client publishing to the nats bus")
	commands := flag.String("commands", "show version", "commands to send tot he device")
	flag.Parse()
	nc, err := nats.Connect(*natsurl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	Cmds := make(map[string]string)
	Cmds["Commands"] = *commands

	jsondata, err := json.Marshal(Cmds)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Send the request
	msg, err := nc.Request(*devicename, jsondata, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Use the response
	log.Printf("Reply: %s", msg.Data)

	// Close the connection
	nc.Close()
}
