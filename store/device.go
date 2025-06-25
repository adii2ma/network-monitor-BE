package store

import (
    "backend/db"
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
    db.RDB.Del(db.Ctx, "device:" + ip) // remove status
    return db.RDB.SRem(db.Ctx, "devices", ip).Err()
}

func GetAllIPs() ([]string, error) {
    return db.RDB.SMembers(db.Ctx, "devices").Result()
}
