package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Ext(event.Name) == ".go" {
					if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("modified file:", event.Name)
						reload()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	reload()

	<-make(chan struct{})
}

var cmd *exec.Cmd

func reload() {
	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			// On Windows, Kill is not implemented, so we need to use taskkill
			// This is a cross-platform way to handle it
			if runtime.GOOS == "windows" {
				exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
			} else {
				log.Println("Error killing process:", err)
			}
		}
	}

	log.Println("Building and restarting application...")

	// The user's requested command sequence
	buildAndRunCmd := "pkill pomodoro-do-ben; killall -9 pomodoro-do-ben; go build . && ./pomodoro-do-ben"
	cmd = exec.Command("sh", "-c", buildAndRunCmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Println("Error starting application:", err)
		return
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			// Don't log "signal: killed" errors, as we are the ones killing it.
			if err.Error() != "signal: killed" {
				log.Println("Application exited with error:", err)
			}
		}
	}()
}
