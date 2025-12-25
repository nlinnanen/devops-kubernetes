package main

import (
	"fmt"
	"net/http"
	"os"
)

func handlePongCreator(count int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count += 1
		fmt.Fprintln(w, "pong ", count)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	count := 0
	http.HandleFunc("/pingpong", handlePongCreator(count))

	fmt.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
