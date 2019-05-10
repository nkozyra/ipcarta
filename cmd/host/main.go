package main

import (
	"flag"

	"github.com/nkozyra/ipcarta"
)

var (
	port          string
	elasticsearch string
)

func init() {
	flag.StringVar(&port, "port", "9999", "Port for host")
	flag.StringVar(&elasticsearch, "el", "", "Elastic search address")
	flag.Parse()
}

func main() {
	ipcarta.Init(ipcarta.Config{
		ElasticSearchHost: elasticsearch,
		Port:              port,
	},
	)
	ipcarta.Serve()
}
