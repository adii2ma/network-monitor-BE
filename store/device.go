package store

import (
    "backend/db"
)

func AddIP(ip string) error {
    return db.RDB.SAdd(db.Ctx, "devices", ip).Err()
}

func DeleteIP(ip string) error {
    db.RDB.Del(db.Ctx, "device:" + ip) // remove status
    return db.RDB.SRem(db.Ctx, "devices", ip).Err()
}

func GetAllIPs() ([]string, error) {
    return db.RDB.SMembers(db.Ctx, "devices").Result()
}
