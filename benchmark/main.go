package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func SetCommandHelper(key, value []byte) []byte {
	buf := []byte{'s', byte(len(key)), byte(len(value))}
	return []byte(fmt.Sprintf("%s%s%s", buf, key, value))
}

func main() {
	var maxKeys int

	flag.IntVar(&maxKeys, "max-keys", 1, "Maximum number of keys to set")
	flag.Parse()

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < maxKeys; i++ {
		conn.Write(
			SetCommandHelper([]byte(fmt.Sprint(i)),
				[]byte(fmt.Sprint(i)),
			),
		)
	}

	conn.Close()
}
