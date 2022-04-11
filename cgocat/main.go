package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <sys/types.h>
// #include <fcntl.h>
// #include <unistd.h>
//
// static int myopen(char *name) {
// 	return open(name, O_RDONLY);
// }
import "C"
import (
	"log"
	"os"
	"unsafe"
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	if len(os.Args) <= 1 {
		cat(0)
	} else {
		for _, fname := range os.Args[1:] {
			csfname := C.CString(fname)
			fd, err := C.myopen(csfname)
			if fd < 0 {
				log.Println(err)
				continue
			}
			err = cat(fd)
			if err != nil {
				log.Println(err)
			}
			C.close(fd)
			C.free(unsafe.Pointer(csfname))
		}
	}
}

func cat(fd C.int) error {
	var buf [C.BUFSIZ]byte
	var n, c C.ssize_t
	var err error
	for {
		n, err = C.read(fd, unsafe.Pointer(&buf), C.size_t(len(buf)))
		if n <= 0 {
			break
		}
		c, err = C.write(1, unsafe.Pointer(&buf), C.size_t(n))
		if c != n {
			break
		}
	}
	return err
}
