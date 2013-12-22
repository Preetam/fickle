package main

import (
	"flag"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listenAddr := flag.String("listen", ":8080", "TCP address to listen on")
	commandLog := flag.String("command-log", "/tmp/fickle.db", "The command log file where received commands are logged")
	debugHTTP := flag.Bool("debug-http", false, "Start an HTTP server for debugging")
	flag.Parse()

	i := NewInstance(*listenAddr, *commandLog)
	if *debugHTTP {
		go StartHttpDebug(i)
	}
	i.Start()
}
