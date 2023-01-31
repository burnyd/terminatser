package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"terminatser/controller/pkg/promexport"

	"github.com/nats-io/nats.go"
)

type GnmiData struct {
	Path   string
	Target string
	Ts     string
	Value  string
}

func printMsg(m *nats.Msg, i int) {
	//log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
	//fmt.Println(m.Subject)
	messages := GnmiData{}

	err := json.Unmarshal([]byte(m.Data), &messages)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	//fmt.Println(messages.Path, messages.Target, messages.Value)
	promexport.ExportToProm(messages.Path, messages.Ts, messages.Target, messages.Value)
}

func main() {
	//init prometheus
	promexport.Init()
	//promexport.RegisterProm()
	natsurl := flag.String("natsurl", "127.0.0.1:4222", "nats url. Default is 127.0.0.1:4222")
	natstopic := flag.String("natstopic", "terminatser", "Topic to send for nats default is terminatser")

	flag.Parse()

	urls := *natsurl
	subj := *natstopic

	opts := []nats.Option{nats.Name(*natstopic)}
	opts = setupConnOptions(opts)

	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	i := 0

	nc.Subscribe(subj, func(msg *nats.Msg) {
		i += 1
		printMsg(msg, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)

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
