package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	mu    sync.Mutex
	count int
)

func handlePong(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	c := count
	mu.Unlock()
	fmt.Fprintln(w, "pong ", c)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	c := count
	mu.Unlock()
	fmt.Fprintln(w, c)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/pingpong", handlePong)
	http.HandleFunc("/count", handleCount)

	fmt.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
