package main

import (
    "fmt"
    "network-monitor/db"
    "network-monitor/pinger"
    "time"
)

func main() {
    db.InitRedis()
    fmt.Println("Starting periodic ping...")

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        <-ticker.C
        pinger.PingAll()
    }
}
