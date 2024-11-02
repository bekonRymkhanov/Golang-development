package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/episode/view", app.episodeView)
	mux.HandleFunc("/episode/add", app.episodeAdd)
	return mux
}
