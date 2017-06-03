package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/julienschmidt/httprouter"
	"github.com/prologic/msgbus"
)

var (
	bind *string
)

func init() {
	flag.String(&bind, "bind", ":8000", "interface and port to bind to")
}

func main() {
	msgbus.Server{bind: bind}.Run()
}
