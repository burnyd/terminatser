package natsreply

import (
	"fmt"
	"log"
	"os"
	"time"

	"terminatser/pkg/eapi"

	"github.com/nats-io/nats.go"
)

type Conn struct {
	Transport string
	Host      string
	Username  string
	Password  string
	Port      int
	Cmds      string
}

func (c Conn) Eapi() string {
	dev := eapi.Conn{
		Transport: c.Transport,
		Host:      c.Host,
		Password:  c.Username,
		Username:  c.Password,
		Port:      80,
		Cmds:      c.Cmds,
	}

	return dev.Connect()
}

/*func (c Conn) Sub(nc nats.Conn, subject string) {
	s, _ := nc.Subscribe(subject, func(msg *nats.Msg) {
		//Need to pass data into eAPI Here in some sort of function?
		//Need to send the response via natsresponse pkg here

		//r := c.Eapi()
		//msg.Respond([]byte(r))
		msg.Respond([]byte("I did it"))

	})

}*/

func (c Conn) StartReply(natsurl string) {
	MySubject, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if MySubject == "" {
		log.Fatal("Empty env var " + MySubject)
	}
	nc, err := nats.Connect(natsurl)
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	sub, _ := nc.Subscribe(MySubject, func(msg *nats.Msg) {
		r := c.Eapi()
		msg.Respond([]byte(r))
	})

	//Need to figure out how to block on this portion.

	time.Sleep(45 * time.Second)

	sub.Unsubscribe()
}

func (c Conn) Init(natsurl string) {
	go func() {
		c.StartReply(natsurl)
	}()
}
