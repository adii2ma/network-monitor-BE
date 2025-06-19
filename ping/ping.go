package pinger

import (
    "fmt"
    "net"
    "time"
    "network-monitor/config"
    "network-monitor/db"
    "network-monitor/models"
)

func PingAll() {
    for _, ip := range config.IPList {
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
