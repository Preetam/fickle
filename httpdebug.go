package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

var m runtime.MemStats

func StartHttpDebug(i *Instance) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		runtime.ReadMemStats(&m)
		err := enc.Encode(m)
		if err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/lexicon", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, i.db.GetRange("\x00", "\xff"))
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
