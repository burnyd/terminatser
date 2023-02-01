package gnmiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/aristanetworks/goarista/gnmi"
	"github.com/nats-io/nats.go"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

type PathOpts struct {
	Paths     []string
	Origin    string
	Timestamp string
	Target    string
	Path      string
	Value     string
	NatsInfo  NatsPublish
	GnmiCfg   gnmi.Config
}

type NatsPublish struct {
	NatsUrls   string
	Subject    string
	ClientName string
}

var cfg = &gnmi.Config{
	Addr:     "172.20.20.2:6030",
	Username: "admin",
	Password: "admin",
}

func (p PathOpts) ClientgNMI() {
	ctx := gnmi.NewContext(context.Background(), cfg)
	client, err := gnmi.Dial(&p.GnmiCfg)
	if err != nil {
		log.Fatal(err)
	}
	subOptions := gnmi.SubscribeOptions{
		Origin: p.Origin,
		Paths:  gnmi.SplitPaths(p.Paths),
		Target: cfg.Addr,
	}
	respChan := make(chan *pb.SubscribeResponse, 128)
	go func() {
		err = gnmi.SubscribeErr(ctx, client, &subOptions, respChan)
		if err != nil {
			log.Fatal(err)
		}
	}()
	for {
		select {
		case response := <-respChan:
			switch resp := response.Response.(type) {
			case *pb.SubscribeResponse_Update:
				t := time.Unix(0, resp.Update.Timestamp).UTC()
				prefix := gnmi.StrPath(resp.Update.Prefix)
				var target string
				if t := resp.Update.Prefix.GetTarget(); t != "" {
					target = "(" + t + ") "
				}
				for _, update := range resp.Update.Update {
					//fmt.Printf("[%s] %sUpdate %s = %s\n",
					//	t.Format(time.RFC3339Nano),
					//	target,
					//	path.Join(prefix, gnmi.StrPath(update.Path)),
					//	gnmi.StrUpdateVal(update),
					//)
					p.Target = target
					p.Timestamp = t.Format(time.RFC3339Nano)
					p.Path = path.Join(prefix, gnmi.StrPath(update.Path))
					p.Value = gnmi.StrUpdateVal(update)
					p.NatsClient(t.Format(time.RFC3339Nano), target, path.Join(prefix, gnmi.StrPath(update.Path)), gnmi.StrUpdateVal(update))
				}
			}
		}
	}
}

func (p PathOpts) NatsClient(ts, target, path, value string) {
	NatsData := map[string]string{
		"ts":     ts,
		"target": target,
		"path":   path,
		"value":  value,
	}

	msg, err := json.Marshal(&NatsData)
	if err != nil {
		fmt.Println("error:", err)
	}

	log.SetFlags(0)

	opts := []nats.Option{nats.Name(p.NatsInfo.ClientName)}

	nc, err := nats.Connect(p.NatsInfo.NatsUrls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Publish(p.NatsInfo.Subject, msg)
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
	}
}
