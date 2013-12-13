package main

import (
	"flag"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listenAddr := flag.String("listen", ":8080", "TCP address to listen on")
	debugHTTP := flag.Bool("debug-http", false, "Start an HTTP server for debugging")
	flag.Parse()

	if *debugHTTP {
		StartHttpDebug()
	}
	i := NewInstance(*listenAddr)
	i.Start()
}
