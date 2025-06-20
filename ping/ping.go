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

    status := models.DeviceStatus{
        IP:       ip,
        Online:   err == nil,
        LastSeen: time.Now().Unix(),
    }

    key := fmt.Sprintf("device:%s", ip)
    db.RDB.HSet(db.Ctx, key, map[string]interface{}{
        "online":    status.Online,
        "last_seen": status.LastSeen,
    })
}
