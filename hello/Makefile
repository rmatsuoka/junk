target = hello
GOOS=linux
GOARCH=amd64
CGO_ENABLE=0

$(target): *.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLE=0 go build -o $@

.PHONY: clean
clean:
	rm -r $(target)

.POHNY: dist
dist: $(target) hello.service install.sh
	tar zcf hello_$(GOOS)_$(GOARCH).tar.gz hello hello.service install.sh
