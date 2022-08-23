package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rmatsuoka/junk/progbar"
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	c := make(chan int)

	go progressbar.Display(2, c)
	for i := 0; i < 100; {
		time.Sleep(time.Millisecond * 125)
		i += rand.Intn(5)
		if i > 100 {
			i = 100
		}
		c <- i
	}
	close(c)
	time.Sleep(time.Millisecond)
}
