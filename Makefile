get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get github.com/gotk3/gotk3/glib
	go get github.com/dlasky/gotk3-layershell/layershell
	go get github.com/joshuarubin/go-sway
	go get github.com/allan-simon/go-singleinstance

build:
	go build -o bin/nwg-bar *.go

install:
	mkdir -p /usr/share/nwg-bar
	cp config/* /usr/share/nwg-bar
	cp bin/nwg-bar /usr/bin

uninstall:
	rm -r /usr/share/nwg-bar
	rm /usr/bin/nwg-bar

run:
	go run *.go
