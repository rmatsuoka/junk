PREFIX=$(HOME)/.local
PROGS= \
	netmon_hls \
	netmon2curl \
	netmon_m3u8_curl \

all:
	@echo run \'make install\'

install:
	mkdir -p $(PREFIX)/bin

	for f in $(PROGS); \
	do \
		rm -f "$(PREFIX)/bin/$$f"; \
		cp "$$f" "$(PREFIX)/bin/$$f"; \
		chmod +x "$(PREFIX)/bin/$$f"; \
	done
