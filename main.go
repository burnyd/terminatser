package main

import (
	"flag"
	"terminatser/pkg/gnmiclient"
	"terminatser/pkg/natsreply"

	"github.com/aristanetworks/goarista/gnmi"
)

func main() {
	natsurl := flag.String("natsurl", "127.0.0.1:4222", "nats url. Default is 127.0.0.1:4222")
	natstopic := flag.String("natstopic", "terminatser", "Topic to send for nats default is terminatser")
	clientname := flag.String("clientname", "eos", "Name of the client publishing to the nats bus")
	gnmitarget := flag.String("gnmitarget", "127.0.0.1:6030", "Address for gNMI")
	gnmiuser := flag.String("gnmiuser", "admin", "username for gnmi")
	gnmipassword := flag.String("gnmipassword", "admin", "password for gnmi")
	gnmipath := flag.String("gnmipath", "/", "path for gnmi")
	ShowCmds := flag.String("showcommands", "show version", "Commands to execute if used")

	flag.Parse()
	//Create gNMI Connection
	s := gnmiclient.PathOpts{
		Paths:  []string{*gnmipath},
		Origin: "openconfig",
		GnmiCfg: gnmi.Config{
			Addr:     *gnmitarget,
			Username: *gnmiuser,
			Password: *gnmipassword,
		},
		NatsInfo: gnmiclient.NatsPublish{
			NatsUrls:   *natsurl,
			Subject:    *natstopic,
			ClientName: *clientname,
		},
	}
	// Start the nats listener
	nr := natsreply.Conn{
		Transport: "http",
		Host:      s.GnmiCfg.Addr,
		Password:  s.GnmiCfg.Username,
		Username:  s.GnmiCfg.Password,
		Port:      80,
		Cmds:      *ShowCmds,
	}
	nr.Init(s.NatsInfo.NatsUrls)
	s.ClientgNMI()
}
