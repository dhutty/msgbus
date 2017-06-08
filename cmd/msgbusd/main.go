package main

import (
	"flag"

	"github.com/prologic/msgbus"
)

var (
	bind string
)

func init() {
	flag.StringVar(&bind, "bind", ":8000", "interface and port to bind to")
}

func main() {
	msgbus.NewServer(bind).ListenAndServe()
}
