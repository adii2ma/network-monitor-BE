package store

import (
	"backend/db"
	"fmt"
	"time"
)

func AddIP(ip string, location string) error {
	err := db.RDB.SAdd(db.Ctx, "devices", ip).Err()
	if err != nil {
		return err
	}
	// Set default status: offline, last_seen = 0, with provided location
	key := "device:" + ip
	fields := []interface{}{
		"online", "false",
		"last_seen", "0",
		"location", location,
	}
	return db.RDB.HMSet(db.Ctx, key, fields...).Err()
}

func DeleteIP(ip string) error {
	db.RDB.Del(db.Ctx, "device:"+ip) // remove status
	return db.RDB.SRem(db.Ctx, "devices", ip).Err()
}

func GetAllIPs() ([]string, error) {
	return db.RDB.SMembers(db.Ctx, "devices").Result()
}

// UpdateDeviceStatus updates the device status and logs changes
func UpdateDeviceStatus(ip string, online bool) error {
	key := "device:" + ip

	// Get current status to check if it changed
	currentStatus, err := db.RDB.HGet(db.Ctx, key, "online").Result()
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	currentOnline := currentStatus == "true"
	timestamp := time.Now().Unix()

	// Update the status
	fields := []interface{}{
		"online", fmt.Sprintf("%t", online),
		"last_seen", fmt.Sprintf("%d", timestamp),
	}

	err = db.RDB.HMSet(db.Ctx, key, fields...).Err()
	if err != nil {
		return err
	}

	// Log status change if it actually changed
	if currentOnline != online {
		status := "offline"
		if online {
			status = "online"
		}

		logEntry := fmt.Sprintf("Device %s went %s at %s",
			ip, status, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))

		// Store in device log
		logKey := "device:log:" + ip
		db.RDB.LPush(db.Ctx, logKey, logEntry)
		db.RDB.LTrim(db.Ctx, logKey, 0, 99) // Keep last 100 entries

		// Also store in global logs
		db.RDB.LPush(db.Ctx, "logs", logEntry)
		db.RDB.LTrim(db.Ctx, "logs", 0, 999) // Keep last 1000 entries
	}

	return nil
}

// GetDeviceLogs returns the logs for a specific device
func GetDeviceLogs(ip string) ([]string, error) {
	logKey := "device:log:" + ip
	return db.RDB.LRange(db.Ctx, logKey, 0, -1).Result()
}
