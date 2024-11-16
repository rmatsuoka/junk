target = hello
GOOS=linux
GOARCH=amd64
CGO_ENABLE=0
disttar=hello_$(GOOS)_$(GOARCH).tar.gz

$(target): *.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLE=0 go build -o $@

.PHONY: clean
clean:
	rm -r $(target)

.POHNY: dist
dist: $(disttar)

$(disttar): $(target) hello.service install.sh
	tar zcf $@ hello hello.service install.sh
