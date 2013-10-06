package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

var m runtime.MemStats

func StartHttpDebug() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		runtime.ReadMemStats(&m)
		err := enc.Encode(m)
		if err != nil {
			log.Println(err)
		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
