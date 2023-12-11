package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joshuarubin/go-sway"
)

func tempDir() string {
	if os.Getenv("TMPDIR") != "" {
		return os.Getenv("TMPDIR")
	} else if os.Getenv("TEMP") != "" {
		return os.Getenv("TEMP")
	} else if os.Getenv("TMP") != "" {
		return os.Getenv("TMP")
	}
	return "/tmp"
}

func readTextFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func configDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return (fmt.Sprintf("%s/nwg-bar", os.Getenv("XDG_CONFIG_HOME")))
	}
	return fmt.Sprintf("%s/.config/nwg-bar", os.Getenv("HOME"))
}

func getDataHome() string {
	if os.Getenv("XDG_DATA_HOME") != "" {
		return os.Getenv("XDG_DATA_HOME")
	}
	return "/usr/share/"
}

func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err == nil {
			fmt.Println("Creating dir:", dir)
		}
	}
}

func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func copyFile(src, dst string) error {
	fmt.Println("Copying file:", dst)

	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func mapOutputs() (map[string]*gdk.Monitor, error) {
	result := make(map[string]*gdk.Monitor)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := sway.New(ctx)
	if err != nil {
		return nil, err
	}

	outputs, err := client.GetOutputs(ctx)
	if err != nil {
		return nil, err
	}

	display, err := gdk.DisplayGetDefault()
	if err != nil {
		return nil, err
	}

	num := display.GetNMonitors()
	for i := 0; i < num; i++ {
		monitor, _ := display.GetMonitor(i)
		geometry := monitor.GetGeometry()
		// assign output to monitor on the basis of the same x, y coordinates
		for _, output := range outputs {
			if int(output.Rect.X) == geometry.GetX() && int(output.Rect.Y) == geometry.GetY() {
				result[output.Name] = monitor
			}
		}
	}
	return result, nil
}

/*
Window on-leave-notify event hides the dock with glib Timeout 500 ms.
We might have left the window by accident, so let's clear the timeout if window re-entered.
Furthermore - hovering a button triggers window on-leave-notify event, and the timeout
needs to be cleared as well.
*/
func cancelClose() {
	if src > 0 {
		glib.SourceRemove(src)
		src = 0
	}
}

func createPixbuf(icon string, size int) (*gdk.Pixbuf, error) {
	if strings.HasPrefix(icon, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(icon, size, size)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return pixbuf, nil
	}

	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		log.Fatal("Couldn't get default theme: ", err)
	}
	pixbuf, err := iconTheme.LoadIcon(icon, size, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		return nil, err
	}
	return pixbuf, nil
}

func launch(command string) {
	// trim % and everything afterwards
	if strings.Contains(command, "%") {
		cutAt := strings.Index(command, "%")
		if cutAt != -1 {
			command = command[:cutAt-1]
		}
	}

	elements := strings.Split(command, " ")

	// find prepended env variables, if any
	envVarsNum := strings.Count(command, "=")
	var envVars []string

	cmdIdx := 0
	lastEnvVarIdx := 0

	if envVarsNum > 0 {
		for idx, item := range elements {
			if strings.Contains(item, "=") {
				lastEnvVarIdx = idx
				envVars = append(envVars, item)
			}
		}
		cmdIdx = lastEnvVarIdx + 1
	}

	cmd := exec.Command(elements[cmdIdx], elements[1+cmdIdx:]...)

	// set env variables
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	msg := fmt.Sprintf("env vars: %s; command: '%s'; args: %s\n", envVars, elements[cmdIdx], elements[1+cmdIdx:])
	println(msg)

	go cmd.Run()

	glib.TimeoutAdd(uint(150), func() bool {
		gtk.MainQuit()
		return false
	})
}
