package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/allan-simon/go-singleinstance"
	/*"github.com/gotk3/gotk3/gtk"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"*/)

const version = "0.0.1"

var (
	configDirectory string
	dataHome        string
	buttons         []Button
)

type Button struct {
	Name string
	Exec string
	Icon string
}

// Flags
var displayVersion = flag.Bool("v", false, "display Version information")

//var cssFileName = flag.String("s", "style.css", "Css file name")
var templateFileName = flag.String("t", "bar.json", "Template file name")

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
				//gtk.MainQuit()
			}
		}
	}()

	// We want the same key/mouse binding to turn the bar off. Kill the running instance and exit.
	lockFilePath := fmt.Sprintf("%s/nwg-bar.lock", tempDir())
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
	// will only be created if does not yet exist
	createDir(configDirectory)

	// Copy default config
	if !pathExists(filepath.Join(configDirectory, "style.css")) {
		res := copyFile(filepath.Join(dataHome, "nwg-bar/style.css"), filepath.Join(configDirectory, "style.css"))
		fmt.Println(res)
	}
	if !pathExists(filepath.Join(configDirectory, "bar.json")) {
		res := copyFile(filepath.Join(dataHome, "nwg-bar/bar.json"), filepath.Join(configDirectory, "bar.json"))
		fmt.Println(res)
	}

	// load JSON template
	p := filepath.Join(configDir(), *templateFileName)
	templateJson, err := readTextFile(p)
	if err != nil {
		log.Fatal(err)
	} else {
		// parse JSON to []Button
		json.Unmarshal([]byte(templateJson), &buttons)
		println(fmt.Sprintf("%v items loaded from %s", len(buttons), p))
	}
}
