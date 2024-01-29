PREFIX ?= /usr
DESTDIR ?= 

get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get github.com/gotk3/gotk3/glib
	go get github.com/dlasky/gotk3-layershell/layershell
	go get github.com/joshuarubin/go-sway
	go get github.com/allan-simon/go-singleinstance

build:
	go build -v -o bin/nwg-bar .

install:
	mkdir -p $(DESTDIR)$(PREFIX)/share/nwg-bar
	cp config/* $(DESTDIR)$(PREFIX)/share/nwg-bar
	mkdir -p $(DESTDIR)$(PREFIX)/share/nwg-bar/images
	cp images/* $(DESTDIR)$(PREFIX)/share/nwg-bar/images
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp bin/nwg-bar $(DESTDIR)$(PREFIX)/bin/nwg-bar

	mkdir -p $(DESTDIR)$(PREFIX)/share/doc/nwg-bar
	cp README.md $(DESTDIR)$(PREFIX)/share/doc/nwg-bar
	mkdir -p $(DESTDIR)$(PREFIX)/share/licenses/nwg-bar
	cp LICENSE $(DESTDIR)$(PREFIX)/share/licenses/nwg-bar

uninstall:
	rm -r $(DESTDIR)$(PREFIX)/share/nwg-bar
	rm $(DESTDIR)$(PREFIX)/bin/nwg-bar

run:
	go run .
