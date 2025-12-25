package main

import (
	"fmt"
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

	countFilePath := os.Getenv("COUNT_PATH")
	if countFilePath == "" {
		countFilePath = "../../pong/count.txt"
	}

	countData, err := os.ReadFile(countFilePath)
	if err != nil {
		http.Error(w, "failed to read count file", http.StatusInternalServerError)
		log.Printf("error reading %s: %v", countFilePath, err)
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
