package main

import (
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/httpRouter"
	"log"
	"net/http"
	"os"
)

func main() {

	label := "Anon"
	if len(os.Args) > 3 {
		label = os.Args[3]
	}
	ipPort := ""
	if len(os.Args) > 2 {
		ipPort = os.Args[1] + ":" + os.Args[2]
	} else if len(os.Args) > 1 {
		ipPort = "localhost:" + os.Args[1]
	} else {
		ipPort = "localhost:6686"
	}
	data.SetDataStore(label, ipPort)

	router := httpRouter.NewRouter()
	if len(os.Args) > 2 {
		log.Fatal(http.ListenAndServe(os.Args[1]+":"+os.Args[2], router))
	} else if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}

}
