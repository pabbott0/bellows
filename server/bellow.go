package main

import (
	"fmt"
	"github.com/SlyMarbo/spdy"
	"github.com/pabbott0/bellows"
	//"net/http"
	"os"
)

func main() {
	b := bellows.Server(os.Args[1])
	conf := b.Context.Conf

	err := spdy.ListenAndServeTLS("localhost:8001", conf.Server.SSLCertPath, conf.Server.SSLKeyPath, b.Handler())
	//err := http.ListenAndServeTLS("localhost:8001", conf.Server.SSLCertPath, conf.Server.SSLKeyPath, b.Handler())

	if err != nil {
		fmt.Println("error: ", err)
	}
}
