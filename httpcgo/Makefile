run: a.out
	./a.out

a.out: get.c get
	cc get.c get

get: get.go
	go build --buildmode=c-shared get.go

clean:
	go clean
	rm -f a.out x
