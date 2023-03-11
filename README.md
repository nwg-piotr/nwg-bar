# nwg-bar

This application is a part of the [nwg-shell](https://nwg-piotr.github.io/nwg-shell) project.

**Contributing:** please read the [general contributing rules for the nwg-shell project](https://nwg-piotr.github.io/nwg-shell/contribution).

nwg-bar is a Golang replacement to the `nwgbar` command (a part of
[nwg-launchers](https://github.com/nwg-piotr/nwg-launchers)), with some improvements. Aimed at sway, works with
wlroots-based compositors only.

The `nwg-bar` command creates a button bar on the basis of a JSON template placed in the `~/.config/nwg-bar/` folder.
By default the command displays a horizontal bar in the center
of the screen. Use command line arguments to change the placement.

![image](https://user-images.githubusercontent.com/20579136/163154930-883140f3-0f69-481f-b07f-cbd4d4c75117.png)


[![Packaging status](https://repology.org/badge/vertical-allrepos/nwg-bar.svg)](https://repology.org/project/nwg-bar/versions)

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

To uninstall run `sudo make uninstall`.

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

*NOTE: for now the `-o` argument works on sway only.*

## Templates

Templates in JSON format should be placed in the `~/.config/nwg-bar` folder. The default `bar.json` template creates
 sample Exit menu for sway on Arch Linux. You may adjust it to your system, and also add as many other templates,
 as you need. Use the `-t somename.json` argument to specify the template name to use.

 ```json
 [
  {
    "label": "Lock",
    "exec": "swaylock -f -c 000000",
    "icon": "/usr/share/nwg-bar/images/system-lock-screen.svg"
  },
  {
    "label": "Logout",
    "exec": "swaymsg exit",
    "icon": "/usr/share/nwg-bar/images/system-log-out.svg"
  },
  {
    "label": "Reboot",
    "exec": "systemctl reboot",
    "icon": "/usr/share/nwg-bar/images/system-reboot.svg"
  },
  {
    "label": "Shutdown",
    "exec": "systemctl -i poweroff",
    "icon": "/usr/share/nwg-bar/images/system-shutdown.svg"
  }
]
 ```

 - `label` field defines the button label;
 - `exec` field defines the command to execute on button click;
 - `icon` field specifies the button icon; you may use a system icon name, like e.g. `system-lock-screen`, or a path to .svg/.png file.

 ## Styling

 Edit the `~/.config/nwg-bar/style.css` file to change the bar appearance. You may also specify another .css file
 (in the same folder) with the `-s somename.css` argument.
