package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	w *fsnotify.Watcher
}

func NewWatcher() *Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return &Watcher{
		w: w,
	}
}

func (w *Watcher) Close() {
	err := w.w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Watcher) process() {
	for {
		select {
		case ev, ok := <-w.w.Events:
			if !ok {
				log.Print("Event channel error:")
				return
			}

			log.Println(ev.String())

			switch {
			case ev.Op&fsnotify.Create != 0:
				if isDirectory(ev.Name) {
					err := w.w.Add(ev.Name)
					if err != nil {
						log.Println("Error add folder into watcher:", err.Error())
					}
				}
			}

		case err, ok := <-w.w.Errors:
			if !ok {
				log.Print("Error channel error:")
				return
			}
			log.Println("Got error:", err.Error())
		}
	}
}

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Println("Check dir error:", err.Error())
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}

func main() {
	log.SetFlags(log.Lshortfile)

	watcher := NewWatcher()
	defer watcher.Close()

	// Run main process
	go watcher.process()

	err := watcher.w.Add("/tmp/mytmp/save")
	if err != nil {
		log.Fatal(err)
	}

	// Wait forever
	done := make(chan bool)
	<-done
}
