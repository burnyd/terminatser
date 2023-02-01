package natsreply

import (
	"fmt"
	"log"
	"os"
	"runtime"
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

func (c Conn) StartReply(natsurl string) {
	MySubject, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if MySubject == "" {
		log.Fatal("Empty env var " + MySubject)
	}
	opts := []nats.Option{nats.Name(MySubject)}
	opts = setupConnOptions(opts)
	nc, err := nats.Connect(natsurl, opts...)
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe(MySubject, func(msg *nats.Msg) {
		r := c.Eapi()
		msg.Respond([]byte(r))
	})

	nc.Flush()

	runtime.Goexit()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
		log.Printf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func (c Conn) Init(natsurl string) {
	go func() {
		c.StartReply(natsurl)
	}()
}
