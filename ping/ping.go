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
     
    // Get existing location to preserve it
    key := fmt.Sprintf("device:%s", ip)
    existingLocation, err := db.RDB.HGet(db.Ctx, key, "location").Result()
    if err != nil {
        // If location doesn't exist, use default
        existingLocation = "Location not set"
    }

    status := models.DeviceStatus{
        IP:       ip,
        Online:   online,
        Location: existingLocation,
        LastSeen: time.Now().Unix(),
    }

    fields := []interface{}{
        "online", fmt.Sprintf("%v", status.Online),
        "last_seen", fmt.Sprintf("%d", status.LastSeen),
        "location", status.Location,
    }
    err = db.RDB.HMSet(db.Ctx, key, fields...).Err()
    if err != nil {
        log.Printf("Failed to update status for %s in Redis: %v\n", ip, err)
    }
}