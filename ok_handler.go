package main

import "net/http"

func OkHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
