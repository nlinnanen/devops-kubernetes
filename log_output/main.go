package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func handleLogCreator(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, getLogStr(id))
	}
}

func getLogStr(id string) string {
	timeStr := time.Now().Format(time.RFC3339)
	return fmt.Sprintf("%s %s", timeStr, id)
}

func logInterval(id string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		log.Println(getLogStr(id))
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	id := uuid.New().String()

	go logInterval(id)

	http.HandleFunc("/log", handleLogCreator(id))

	fmt.Println("Server started in port " + port)
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
