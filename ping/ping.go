package pinger

import (
    "fmt"
    
    "time"
    "log"
    "github.com/go-ping/ping"
    "backend/db"
    "backend/models"
	"backend/store"
)

func PingAll() {
    fmt.Println("PingAll called at", time.Now())
    ips, err := store.GetAllIPs()
    if err != nil {
        log.Println("Failed to get IPs from store:", err)
        return
    }

    for _, ip := range ips {
        go pingDevice(ip)
    }
}

func pingDevice(ip string) {
    pinger, err := ping.NewPinger(ip)
    if err != nil {
        log.Printf("Failed to create pinger for %s: %v\n", ip, err)
        return
    }

    pinger.Count = 1
    pinger.Timeout = time.Second
    pinger.SetPrivileged(true) // May require root/admin on some OSes

    err = pinger.Run()
    if err != nil {
        log.Printf("Ping error for %s: %v\n", ip, err)
    }

    stats := pinger.Statistics()
    online := stats.PacketsRecv > 0

    if online {
        fmt.Printf("[✓] %s is online\n", ip)
    } else {
        fmt.Printf("[✗] %s is offline\n", ip)
    }

    status := models.DeviceStatus{
        IP:       ip,
        Online:   online,
        LastSeen: time.Now().Unix(),
    }

    key := fmt.Sprintf("device:%s", ip)
    db.RDB.HSet(db.Ctx, key, "online", status.Online, "last_seen", status.LastSeen)
}