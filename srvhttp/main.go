package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8000", "address:port")
)

func main() {
	flag.Parse()

	dir := "."
	if flag.NArg() != 0 {
		dir = flag.Arg(0)
	}
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir(dir))))
}