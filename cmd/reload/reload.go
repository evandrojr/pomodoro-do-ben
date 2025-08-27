
package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					reload()
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
			log.Println("Error killing process:", err)
		}
	}

	log.Println("Reloading...")
	cmd = exec.Command("go", "run", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Println("Error starting application:", err)
		return
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Println("Application exited with error:", err)
		}
	}()
}
