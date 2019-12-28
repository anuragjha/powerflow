package main

import (
	"github.com/edgexfoundry/powerflow/registerPowerflow/data"
	"github.com/edgexfoundry/powerflow/registerPowerflow/routing"
	"log"
	"net/http"
	"os"
)

func main() { // for register service

	ipPort := ""
	if len(os.Args) > 2 {
		ipPort = os.Args[1] + ":" + os.Args[2]
	} else if len(os.Args) > 1 {
		ipPort = "127.0.0.1:" + os.Args[1]
	} else {
		ipPort = "127.0.0.1:6680"
	}
	data.SetDataStore(ipPort)

	router := routing.NewRouter()
	if len(os.Args) > 2 {
		log.Fatal(http.ListenAndServe(os.Args[1]+":"+os.Args[2], router))
	} else if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6680", router))
	}

}
