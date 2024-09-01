package main

import "C"
import (
	"errors"
	"io"
	"log"
	"net/http"
	"unsafe"
)

//export get
func get(url *C.char, n C.int, buf *C.char) C.int {
	gourl := C.GoString(url)
	gobuf := unsafe.Slice((*byte)(unsafe.Pointer(buf)), n)

	println("http.get")
	res, err := http.Get(gourl)
	if err != nil {
		log.Print(err)
		return -1
	}
	defer res.Body.Close()

	println("io.readfull")
	size, err := io.ReadFull(res.Body, gobuf)
	if errors.Is(err, io.ErrUnexpectedEOF) {
		err = nil
	}
	if err != nil {
		log.Print(err)
		return -1
	}
	return C.int(size)
}

func main() {}
