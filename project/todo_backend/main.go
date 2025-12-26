package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	todos = []string{}
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		r.Header.Set("Content-Type", "application/json")
		body, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, "failed to marshal todos", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
		return
	} else if r.Method == http.MethodPost {
		todo := r.FormValue("todo")
		if todo == "" {
			http.Error(w, "todo cannot be empty", http.StatusBadRequest)
			return
		}
		todos = append(todos, todo)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "todo added")
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/todos", handler)

	log.Println("Server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
