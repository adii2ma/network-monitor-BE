package pinger

import (
	"backend/db"
	"backend/store"
	"fmt"
	"log"
	"time"

	"github.com/go-ping/ping"
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

	// Use the new UpdateDeviceStatus function which handles logging
	err = store.UpdateDeviceStatus(ip, online)
	if err != nil {
		log.Printf("Failed to update status for %s: %v\n", ip, err)
		return
	}

	// Update location separately if needed (since UpdateDeviceStatus doesn't handle location)
	err = db.RDB.HMSet(db.Ctx, key, "location", existingLocation).Err()
	if err != nil {
		log.Printf("Failed to update location for %s: %v\n", ip, err)
	}
}
