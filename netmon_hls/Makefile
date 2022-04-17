PREFIX=$(HOME)/.local

netmon_hls.fixed: netmon_hls
	sed 's:^libpath=.*###Makefile###:libpath=$(PREFIX)/lib/netmon_hls:' netmon_hls > $@

install: netmon_hls.fixed
	mkdir -p $(PREFIX)/bin
	mkdir -p $(PREFIX)/lib/netmon_hls

	rm -f $(PREFIX)/bin/netmon_hls
	cp netmon_hls.fixed $(PREFIX)/bin/netmon_hls
	chmod +x $(PREFIX)/bin/netmon_hls

	rm -rf $(PREFIX)/lib/netmon_hls
	cp -r lib $(PREFIX)/lib/netmon_hls
uninstall:
	rm -f $(PREFIX)/bin/netmon_hls
	rm -rf $(PREFIX)/lib/netmon_hls
clean:
	rm -f netmon_hls.fixed
