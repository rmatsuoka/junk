package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

var (
	addr     = flag.String("addr", ":8000", "`address:port`")
	mflag    = flag.Bool("m", false, "handle multipart/form-data POST request")
	pathflag = flag.String("path", "/upload", "`path` to handle POST request")
)

func Usage() {
	fmt.Fprintln(os.Stderr, "usage: postsrv [-flags] cmd [arg ...]")
	flag.PrintDefaults()
}

func postHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	c := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	if *mflag {
		r, _ := req.MultipartReader()
		c.Stdin = NewReader(r)
	} else {
		c.Stdin = req.Body
	}
	c.Stdout = w
	c.Stderr = os.Stderr

	c.Run()
	req.Body.Close()
}

type Reader struct {
	r   *multipart.Reader
	p   *multipart.Part
	err error
}

func NewReader(r *multipart.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.p == nil || r.err == io.EOF {
		r.p, r.err = r.r.NextPart()
	}
	if r.err != nil {
		return 0, r.err
	}
	n, r.err = r.p.Read(p)
	if r.err == io.EOF {
		r.p.Close()
		return n, nil
	}
	return n, r.err
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if flag.NArg() == 0 {
		Usage()
		os.Exit(1)
	}

	if _, err := exec.LookPath(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc(*pathflag, postHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
