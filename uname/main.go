package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

var (
	a = flag.Bool("a", false, "all")
	s = flag.Bool("s", false, "system name")
	n = flag.Bool("n", false, "nodename")
	r = flag.Bool("r", false, "release")
	v = flag.Bool("v", false, "version")
	m = flag.Bool("m", false, "machine")
	d = flag.Bool("d", false, "domain name")
)

func main() {
	flag.Parse()

	var u unix.Utsname
	if err := unix.Uname(&u); err != nil {
		log.Fatal(err)
	}

	x := []struct {
		flag *bool
		data *[65]byte
	}{
		{s, &u.Sysname},
		{n, &u.Nodename},
		{r, &u.Release},
		{v, &u.Version},
		{m, &u.Machine},
		{d, &u.Domainname},
	}

	spc := ""
	for _, s := range x {
		if *s.flag || *a {
			fmt.Printf("%s%s", spc, *s.data)
			spc = " "
		}
	}

	// if print nothing
	if spc == "" {
		fmt.Printf("%s", u.Sysname)
	}
	fmt.Println()
}
