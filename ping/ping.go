package pinger

import (
    "fmt"
    "net"
    "time"
    "log"
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
    timeout := 1 * time.Second
    _, err := net.DialTimeout("ip4:icmp", ip, timeout)
    online :=err == nil 
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
    db.RDB.HSet(db.Ctx, key, map[string]interface{}{
        "online":    status.Online,
        "last_seen": status.LastSeen,
    })
}
