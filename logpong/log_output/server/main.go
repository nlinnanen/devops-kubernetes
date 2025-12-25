package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		http.Error(w, "FILE_PATH not set", http.StatusInternalServerError)
		log.Printf("error: FILE_PATH not set")
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		log.Printf("error reading %s: %v", filePath, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/log", handler)

	fmt.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
