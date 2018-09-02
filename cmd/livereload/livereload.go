package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/omeid/go-livereload"
	"github.com/omeid/log"
)

var (
	livereloadAddr = flag.String("livereload", ":35729", "livereload servera addr.")
	serverAddr     = flag.String("server", ":8082", "static server addr. Requires -serve ")
	serve          = flag.String("serve", "", "static files folders.")
	strip          = flag.String("strip", "", "path to strip from static files.")
)

func main() {
	flag.Parse()
	log := log.New("livereload ")

	globs := flag.Args()
	if len(globs) == 0 {
		log.Fatalf("No Files To Watch.")
	}

	var files []string

	for _, glob := range globs {

		matches, err := filepath.Glob(glob)
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, matches...)
	}

	if len(files) == 0 {
		log.Fatalf("No Files Matched The Globs.")
	}

	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watch.Close()

	for _, file := range files {
		err = watch.Add(file)
		if err != nil {
			log.Error(err)
			return
		}
	}

	if *serve != "" {
		go func() {
			static := http.StripPrefix(*strip, http.FileServer(http.Dir(*serve)))
			log.Infof("Serving '%s' at %s", *serve, *serverAddr)
			log.Fatal(http.ListenAndServe(*serverAddr, static))
		}()
	}
	lr := livereload.New("go-livereload")
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/livereload.js", livereload.LivereloadScript)
		mux.Handle("/", lr)
		log.Infof("Serving livereload at %s", *livereloadAddr)
		log.Fatal(http.ListenAndServe(*livereloadAddr, mux))
	}()

	for {
		select {
		case event := <-watch.Events:
			if event.Op&(fsnotify.Rename|fsnotify.Create|fsnotify.Write) > 0 {
				log.Infof("Reloading %s", event.Name)
				lr.Reload(event.Name, true)
			}
		case err := <-watch.Errors:
			if err != nil {
				log.Error(err)
			}
		}
	}
}
