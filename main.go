package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting server...")

	if len(os.Args) != 2 {
		slog.Error("Invalid number of arguments")
		os.Exit(1)
	}

	jc, err := NewJournalConfig(os.Args[1])
	if err != nil {
		slog.Error("Failed to load journal config")
		os.Exit(1)
	}

	handler := HandleMemosBuilder(jc)

	mux := http.NewServeMux()
	mux.HandleFunc("/memos", handler)
	mux.HandleFunc("/ok", OkHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Failed to start server")
		os.Exit(1)
	}

	slog.Info("Server exited successfully")
}
