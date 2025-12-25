package main

import (
	"fmt"
	"net/http"
	"os"
)

func handlePongCreator(count int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count += 1

		err := writeCountToFile(count)
		if err != nil {
			http.Error(w, "failed to write count to file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "pong ", count)
	}
}

func writeCountToFile(count int) error {
	path := os.Getenv("COUNT_FILE")
	if path == "" {
		path = "./count.txt"
	}

	data := fmt.Sprintf("%d", count)
	return os.WriteFile(path, []byte(data), 0644)
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
