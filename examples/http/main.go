package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/blang/mpv"
)

func main() {
	mpvll := mpv.NewIPCClient("/tmp/mpvsocket")
	mpvc := mpv.Client{mpvll}

	// Build in low level json api
	http.Handle("/lowlevel", mpv.HTTPServerHandler(mpvll))

	// Your own api based on mpv.Client
	http.HandleFunc("/file/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s", r.RequestURI)
		path := strings.Replace(r.RequestURI, "/file/", "", 1)
		fmt.Fprintln(w, path)
		log.Println(mpvc.Loadfile(path, mpv.LoadFileModeAppendPlay))
	}))
	http.HandleFunc("/playlist/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s", r.RequestURI)
		path := strings.Replace(r.RequestURI, "/playlist/", "", 1)
		fmt.Fprintln(w, path)
		log.Println(mpvc.LoadList(path, mpv.LoadListModeReplace))
	}))
	http.HandleFunc("/cmd/seek", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mpvc.Seek(10, mpv.SeekModeAbsolute))
	}))
	http.HandleFunc("/cmd/pause", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mpvc.SetPause(true))
	}))
	http.HandleFunc("/cmd/prev", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mpvc.PlaylistPrevious())
	}))
	http.HandleFunc("/cmd/next", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mpvc.PlaylistNext())
	}))
	http.HandleFunc("/value/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.Replace(r.RequestURI, "/value/", "", 1)
		value, err := mpvc.GetProperty(name)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprintf(w, "Value: %s", value)
	}))
	http.HandleFunc("/cmd/fullscreen", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := mpvc.SetProperty("fullscreen", true)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		fmt.Fprintf(w, "ok")
	}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
