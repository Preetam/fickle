package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listenAddr := flag.String("listen", ":8080", "TCP address to listen on")
	commandLog := flag.String("command-log", "/tmp/fickle.db", "The command log file where received commands are logged")
	debugHTTP := flag.Bool("debug-http", false, "Start an HTTP server for debugging")
	replicas := flag.String("replicas", "", "A comma-separated addresses of replicas")
	flag.Parse()

	i := NewInstance(*listenAddr, *commandLog)
	if *debugHTTP {
		go StartHttpDebug(i)
	}

	for _, address := range strings.Split(*replicas, ",") {
		if address != "" {
			log.Println("Added replica:", address)
			i.AddReplica(address)
		}
	}

	i.Start()
}
