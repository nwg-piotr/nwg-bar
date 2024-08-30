package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/allan-simon/go-singleinstance"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.1.6"

var (
	configDirectory  string
	dataHome         string
	buttons          []Button
	src              glib.SourceHandle
	outerOrientation gtk.Orientation
	innerOrientation gtk.Orientation
)

type Button struct {
	Icon  string
	Label string
	Exec  string
}

// Flags
var alignment = flag.String("a", "middle", "Alignment in full width/height: \"start\" or \"end\"")
var full = flag.Bool("f", false, "take Full screen width/height")
var imgSize = flag.Int("i", 48, "Icon size")
var targetOutput = flag.String("o", "", "name of Output to display the bar on")
var position = flag.String("p", "center", "Position: \"bottom\", \"top\", \"left\" or \"right\"")

var marginTop = flag.Int("mt", 0, "Margin Top")
var marginLeft = flag.Int("ml", 0, "Margin Left")
var marginRight = flag.Int("mr", 0, "Margin Right")
var marginBottom = flag.Int("mb", 0, "Margin Bottom")

var cssFileName = flag.String("s", "style.css", "csS file name")
var templateFileName = flag.String("t", "bar.json", "Template file name")
var displayVersion = flag.Bool("v", false, "display Version information")
var exclusiveZone = flag.Bool("x", false, "open on top layer witch eXclusive zone")
var gtkTheme = flag.String("g", "", "GTK theme name")

func main() {
	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwg-bar version %s\n", version)
		os.Exit(0)
	}

	// Gentle SIGTERM handler thanks to reiki4040 https://gist.github.com/reiki4040/be3705f307d3cd136e85
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM {
				println("SIGTERM received, bye bye!")
				gtk.MainQuit()
			}
		}
	}()

	// We want the same key/mouse binding to turn the bar off. Kill the running instance and exit.
	currentUserId := "no-user"
	currentUser, err := user.Current()
	if err == nil {
		currentUserId = currentUser.Uid
	}
	lockFilePath := fmt.Sprintf("%s/%s-nwg-bar.lock", tempDir(), currentUserId)
	lockFile, err := singleinstance.CreateLockFile(lockFilePath)
	if err != nil {
		pid, err := readTextFile(lockFilePath)
		if err == nil {
			i, err := strconv.Atoi(pid)
			if err == nil {
				syscall.Kill(i, syscall.SIGTERM)
			}
		}
		os.Exit(0)
	}
	defer lockFile.Close()

	dataHome = getDataHome()

	configDirectory = configDir()
	// will only be created if it does not yet exist
	createDir(configDirectory)

	// Copy default config
	if !pathExists(filepath.Join(configDirectory, "style.css")) {
		err := copyFile(filepath.Join(dataHome, "nwg-bar/style.css"), filepath.Join(configDirectory, "style.css"))
		if err != nil {
			log.Fatal("error copying default style.css: ", err)
			return
		}
	}
	if !pathExists(filepath.Join(configDirectory, "bar.json")) {
		err := copyFile(filepath.Join(dataHome, "nwg-bar/bar.json"), filepath.Join(configDirectory, "bar.json"))
		if err != nil {
			log.Fatal("error copying default bar.json: ", err)
			return
		}
	}

	// load JSON template
	if !strings.HasPrefix(*templateFileName, "/") {
		*templateFileName = filepath.Join(configDirectory, *templateFileName)
	}
	templateJson, err := readTextFile(*templateFileName)
	if err != nil {
		log.Fatal(err)
	} else {
		// parse JSON to []Button
		err := json.Unmarshal([]byte(templateJson), &buttons)
		if err != nil {
			log.Fatal("error parsing config: ", err)
			return
		} else {
			println(fmt.Sprintf("%v items loaded from template %s", len(buttons), *templateFileName))
		}
	}

	// load style sheet
	if !strings.HasPrefix(*cssFileName, "/") {
		*cssFileName = filepath.Join(configDirectory, *cssFileName)
	}

	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()
	if *gtkTheme != "" {
		err = settings.SetProperty("gtk-theme-name", *gtkTheme)
		if err != nil {
			fmt.Printf("Unable to set theme: %s\n", err)
		} else {
			fmt.Printf("User demanded theme: %s\n", *gtkTheme)
		}
	}

	screen, _ := gdk.ScreenGetDefault()

	cssProvider, _ := gtk.CssProviderNew()

	err = cssProvider.LoadFromPath(*cssFileName)
	if err != nil {
		fmt.Printf("%s file erroneous or not found, using GTK styling\n", *cssFileName)
	} else {
		fmt.Printf("Using style: %s\n", *cssFileName)
		gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	visual, _ := screen.GetRGBAVisual()
	if visual != nil && screen.IsComposited() {
		win.SetVisual(visual)
	}

	layershell.InitForWindow(win)

	// if -o argument given
	var output2mon map[string]*gdk.Monitor
	if *targetOutput != "" {
		// We want to assign layershell to a monitor, but we only know the output name!
		output2mon, err = mapOutputs()
		if err == nil {
			layershell.SetMonitor(win, output2mon[*targetOutput])
		} else {
			fmt.Println(err)
		}
	}

	outerOrientation = gtk.ORIENTATION_VERTICAL
	innerOrientation = gtk.ORIENTATION_HORIZONTAL

	if *position == "bottom" || *position == "top" {
		if *position == "bottom" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
		} else {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, true)
		}

		outerOrientation = gtk.ORIENTATION_VERTICAL
		innerOrientation = gtk.ORIENTATION_HORIZONTAL

		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, *full)
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, *full)
	}

	if *position == "left" || *position == "right" {
		if *position == "left" {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
		} else {
			layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)
		}

		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, *full)
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, *full)

		outerOrientation = gtk.ORIENTATION_HORIZONTAL
		innerOrientation = gtk.ORIENTATION_VERTICAL
	}

	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, *marginTop)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, *marginLeft)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, *marginRight)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_BOTTOM, *marginBottom)

	if !*exclusiveZone {
		layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_OVERLAY)
		layershell.SetExclusiveZone(win, -1)
	} else {
		layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)
		layershell.SetExclusiveZone(win, 0)
	}

	layershell.SetKeyboardMode(win, layershell.LAYER_SHELL_KEYBOARD_MODE_EXCLUSIVE)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Close the window on leave, but not immediately, to avoid accidental closes
	win.Connect("leave-notify-event", func() {
		src = glib.TimeoutAdd(uint(500), func() bool {
			gtk.MainQuit()
			src = 0
			return false
		})
	})

	win.Connect("enter-notify-event", func() {
		cancelClose()
	})

	win.Connect("key-release-event", func(window *gtk.Window, event *gdk.Event) {
		key := &gdk.EventKey{Event: event}
		if key.KeyVal() == gdk.KEY_Escape {
			gtk.MainQuit()
		}
	})

	outerBox, _ := gtk.BoxNew(outerOrientation, 0)
	outerBox.SetProperty("name", "outer-box")
	win.Add(outerBox)

	alignmentBox, _ := gtk.BoxNew(innerOrientation, 0)
	outerBox.PackStart(alignmentBox, true, false, 0)

	mainBox, _ := gtk.BoxNew(innerOrientation, 0)
	mainBox.SetHomogeneous(true)
	mainBox.SetProperty("name", "inner-box")

	if *alignment == "start" {
		alignmentBox.PackStart(mainBox, false, true, 0)
	} else if *alignment == "end" {
		alignmentBox.PackEnd(mainBox, false, true, 0)
	} else {
		alignmentBox.PackStart(mainBox, true, false, 0)
	}

	for _, b := range buttons {
		button, _ := gtk.ButtonNew()
		button.SetProperty("use-underline", true)

		if b.Icon != "" {
			button.SetAlwaysShowImage(true)
			button.SetImagePosition(gtk.POS_TOP)

			pixbuf, err := createPixbuf(b.Icon, *imgSize)
			var img *gtk.Image
			if err == nil {
				img, _ = gtk.ImageNewFromPixbuf(pixbuf)
			} else {
				img, _ = gtk.ImageNewFromIconName("image-missing", gtk.ICON_SIZE_INVALID)
			}
			button.SetImage(img)
		}

		if b.Label != "" {
			button.SetLabel(b.Label)
		}

		button.Connect("enter-notify-event", func() {
			cancelClose()
		})

		exec := b.Exec

		button.Connect("clicked", func() {
			launch(exec)
		})

		mainBox.PackStart(button, true, true, 0)
	}

	win.ShowAll()
	gtk.Main()
}
