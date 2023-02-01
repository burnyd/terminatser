package eapi

import (
	"fmt"
	"strings"

	"github.com/aristanetworks/goeapi"
)

type Conn struct {
	Transport string
	Host      string
	Username  string
	Password  string
	Port      int
	Cmds      string
}

func (c *Conn) Connect() string {
	connect, err := goeapi.Connect(c.Transport, strings.ReplaceAll(c.Host, ":6030", ""), c.Username, c.Password, c.Port)
	if err != nil {
		fmt.Println(err)
	}

	runcommands, err := connect.RunCommands([]string{c.Cmds}, "text")
	if err != nil {
		fmt.Println(err)
	}

	s := fmt.Sprintf("%v", runcommands.Result)
	return s
}
