package store

import (
	"backend/db"
	"backend/mail"
	"fmt"
	"log"
	"time"
)

func AddIP(ip string, location string, name string) error {
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
		"name", name,
	}
	log.Printf("Adding device %s with location %s and name %s", ip, location, name)
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
	location, err := db.RDB.HGet(db.Ctx, key, "location").Result()
	if err != nil {
		log.Printf("Error fetching location for IP %s: %v", ip, err)

	}
	fmt.Printf("Device %s is located at: %s\n", ip, location)

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

		// Send email notification
		emailSubject := fmt.Sprintf("Network Alert: Device %s Status Change", ip)
		emailMessage := fmt.Sprintf("Device <strong>%s</strong> has gone <strong>%s</strong> at <strong>%s</strong>",
			ip, status, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))

		// Send email in a goroutine to avoid blocking
		go func() {
			err := mail.SendNotificationMail(emailSubject, emailMessage)
			if err != nil {
				log.Printf("Failed to send notification email for device %s: %v", ip, err)
			}
		}()
	}

	return nil
}

// GetDeviceLogs returns the logs for a specific device
func GetDeviceLogs(ip string) ([]string, error) {
	logKey := "device:log:" + ip
	return db.RDB.LRange(db.Ctx, logKey, 0, -1).Result()
}
