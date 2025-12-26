package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /log")

	logFilePath := os.Getenv("LOG_PATH")
	if logFilePath == "" {
		http.Error(w, "LOG_PATH not set", http.StatusInternalServerError)
		log.Printf("error: LOG_PATH not set")
		return
	}

	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		log.Printf("error reading %s: %v", logFilePath, err)
		return
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		http.Error(w, "FILE_PATH not set", http.StatusInternalServerError)
		log.Printf("error: FILE_PATH not set")
		return
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		log.Printf("error reading %s: %v", filePath, err)
		return
	}

	message := os.Getenv("MESSAGE")
	if message == "" {
		http.Error(w, "MESSAGE not set", http.StatusInternalServerError)
		log.Printf("error: MESSAGE not set")
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

	data := append([]byte("file content:\n"), fileData...)
	data = append(data, []byte("env variable: MESSAGE="+message+"\n")...)
	data = append(data, []byte("\n")...)
	data = append(data, logData...)
	data = append(data, []byte("Ping / Pongs: ")...)
	data = append(data, countData...)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /log")
	fmt.Fprintln(w, "Log endpoint is working")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8880"
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/log", logHandler)

	log.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
