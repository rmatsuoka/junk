run: a.out
	./a.out

a.out: get.c x
	cc get.c x

x: add.go
	go build --buildmode=c-shared

clean:
	go clean
	rm -f a.out x
