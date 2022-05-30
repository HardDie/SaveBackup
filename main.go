package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	w *fsnotify.Watcher
	b Backup
}

func NewWatcher() *Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return &Watcher{
		w: w,
		b: NewBackup(GetConfig().Track, GetConfig().Name),
	}
}
func (w *Watcher) Close() {
	err := w.w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func (w *Watcher) Process() {
	previousRename := ""
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
				if previousRename != "" {
					err := w.b.Rename(previousRename, ev.Name)
					previousRename = ""
					if err != nil {
						log.Println("Error backup file:", err.Error())
					}
				} else {
					err := w.b.Create(ev.Name)
					if err != nil {
						log.Println("Error backup file:", err.Error())
					}
				}
				if IsDirectory(ev.Name) {
					err := w.w.Add(ev.Name)
					if err != nil {
						log.Println("Error add folder into watcher:", err.Error())
					}
				}
			case ev.Op&fsnotify.Write != 0:
				err := w.b.Changes(ev.Name)
				if err != nil {
					log.Println("Error backup file:", err.Error())
				}
			case ev.Op&fsnotify.Remove != 0:
				err := w.b.Delete(ev.Name)
				if err != nil {
					log.Println("Error backup file:", err.Error())
				}
			case ev.Op&fsnotify.Rename != 0:
				previousRename = ev.Name
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

func main() {
	log.SetFlags(log.Lshortfile)

	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	watcher := NewWatcher()
	defer watcher.Close()

	err = watcher.b.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.w.Add(GetConfig().Track)
	if err != nil {
		log.Fatal(err)
	}

	// Run main process
	go watcher.Process()

	// Wait forever
	done := make(chan bool)
	<-done
}
