package store

import (
    "backend/db"
)

func AddIP(ip string) error {
    err := db.RDB.SAdd(db.Ctx, "devices", ip).Err()
    if err != nil {
        return err
    }
    // Set default status: offline, last_seen = 0
    key := "device:" + ip
    return db.RDB.HSet(db.Ctx, key, "online", "false", "last_seen", "0").Err()
}

func DeleteIP(ip string) error {
    db.RDB.Del(db.Ctx, "device:" + ip) // remove status
    return db.RDB.SRem(db.Ctx, "devices", ip).Err()
}

func GetAllIPs() ([]string, error) {
    return db.RDB.SMembers(db.Ctx, "devices").Result()
}
