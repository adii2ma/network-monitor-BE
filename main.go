package main

import (
    "fmt"
    "net/http"
    "backend/api"
    "backend/db"
    "backend/ping"
    "time"
)

func main() {
    db.InitRedis()
    fmt.Println("Server started at :8080")

    // Start periodic ping in a goroutine
    go func() {
        fmt.Println("Starting periodic ping...")
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()

        for {
            <-ticker.C
            pinger.PingAll()
        }
    }()

    // Start HTTP server (blocks main thread)
    http.HandleFunc("/add", api.AddHandler)
    http.HandleFunc("/delete", api.DeleteHandler)
    http.HandleFunc("/status", api.StatusHandler)

    // If ListenAndServe fails, log the error
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Server error:", err)
    }
}

