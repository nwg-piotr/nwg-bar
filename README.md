# nwg-bar

Golang replacement to the `nwgbar` command (a part of [nwg-launchers](https://github.com/nwg-piotr/nwg-launchers)),
with some improvements. Aimed at sway, works with wlroots-based compositors only.

## Installation

### Requirements

- `go` 1.16 (just to build)
- `gtk3`
- `gtk-layer-shell`

### Steps

1. Clone the repository, cd into it.
2. Install golang libraries with `make get`. First time it may take ages, be patient.
3. `make build`
4. `sudo make install`

If your machine is x86_64, you may skip 2 and 3, and just install the provided binary with `sudo make install`.

## Running

```text
Usage of nwg-bar:
  -a string
    	Alignment in full width/height: "start" or "end" (default "middle")
  -f	take Full screen width/height
  -i int
    	Icon size (default 48)
  -mb int
    	Margin Bottom
  -ml int
    	Margin Left
  -mr int
    	Margin Right
  -mt int
    	Margin Top
  -o string
    	name of Output to display the bar on
  -p string
    	Position: "bottom", "top", "left" or "right" (default "center")
  -s string
    	csS file name (default "style.css")
  -t string
    	Template file name (default "bar.json")
  -v	display Version information
  -x	open on top layer witch eXclusive zone
```
