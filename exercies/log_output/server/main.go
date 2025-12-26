package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := os.Getenv("LOG_PATH")
	if filePath == "" {
		http.Error(w, "LOG_PATH not set", http.StatusInternalServerError)
		log.Printf("error: LOG_PATH not set")
		return
	}

	logData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		log.Printf("error reading %s: %v", filePath, err)
		return
	}

	countURL := os.Getenv("COUNT_URL")
	if countURL == "" {
		countURL = "http://localhost:8080/pingpong/count"
	}

	resp, err := http.Get(countURL)
	if err != nil {
		http.Error(w, "failed to fetch count", http.StatusInternalServerError)
		log.Printf("error fetching %s: %v", countURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "failed to fetch count", http.StatusInternalServerError)
		log.Printf("unexpected status from %s: %s", countURL, resp.Status)
		return
	}

	countData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read count response", http.StatusInternalServerError)
		log.Printf("error reading response from %s: %v", countURL, err)
		return
	}

	data := append(logData, []byte("\nPing / Pongs: ")...)
	data = append(data, countData...)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8880"
	}
	http.HandleFunc("/log", handler)

	fmt.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
