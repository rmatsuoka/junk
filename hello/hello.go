package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

var (
	addr     = flag.String("addr", ":8080", "addr:port")
	certFile = flag.String("cert", "", "TLS certFile")
	keyFile  = flag.String("key", "", "TLS keyFile")
)

func main() {
	flag.Parse()
	enableTLS := *certFile != "" || *keyFile != ""

	if enableTLS {
		if *certFile == "" {
			log.Fatal("no certFile")
		}
		if *keyFile == "" {
			log.Fatal("no keyFile")
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		slog.InfoContext(req.Context(), "request",
			"method", req.Method,
			"host", req.Host,
			"url", req.URL.String(),
			"pattern", req.Pattern,
			"remoteAddr", req.RemoteAddr,
			"userAgent", req.UserAgent(),
			"tls", req.TLS,
		)
		fmt.Fprintf(w, "hello, %s\n", req.RemoteAddr)
	})

	var err error
	if enableTLS {
		err = http.ListenAndServeTLS(*addr, *certFile, *keyFile, nil)
	} else {
		err = http.ListenAndServe(*addr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
