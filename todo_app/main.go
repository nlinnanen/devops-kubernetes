package main

import (
    "fmt"
    "net/http"
    "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, world")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/", handler)

    fmt.Println("Server started in port " + port)
    addr := ":" + port
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        panic(err)
    }
}

