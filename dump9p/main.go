package main

// dump9p -- convert 9P2000 data into text.

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"9fans.net/go/plan9"
)

var (
	stdout   io.Writer = os.Stdout
	exitCode           = 0
	uflag              = flag.Bool("u", false, "unbuffer")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("dump9p: ")
	flag.Parse()

	if !*uflag {
		stdout = bufio.NewWriter(os.Stdout)
	}

	if flag.NArg() == 0 {
		look9p(os.Stdin)
	} else {
		for _, fname := range flag.Args() {
			f, err := os.Open(fname)
			if err != nil {
				log.Print(err)
				exitCode = 1
				continue
			}
			look9p(f)
			f.Close()
		}
	}

	if w, ok := stdout.(*bufio.Writer); ok {
		w.Flush()
	}
}

func look9p(f *os.File) {
	name := f.Name()
	var r io.Reader = f
	if !*uflag {
		r = bufio.NewReader(r)
	}
	for {
		f, err := plan9.ReadFcall(r)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("readfcall %s: %v", name, err)
			return
		}
		fmt.Fprintln(stdout, f)
	}
}
