package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "backend/store"
    "backend/db"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
    ip := r.URL.Query().Get("ip")
    if ip == "" {
        http.Error(w, "Missing ?ip= parameter", http.StatusBadRequest)
        return
    }
    if err := store.AddIP(ip); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Added IP: %s\n", ip)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
    ip := r.URL.Query().Get("ip")
    if ip == "" {
        http.Error(w, "Missing ?ip= parameter", http.StatusBadRequest)
        return
    }
    if err := store.DeleteIP(ip); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Deleted IP: %s\n", ip)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    ips, err := store.GetAllIPs()
    if err != nil {
        http.Error(w, "Failed to get IP list", http.StatusInternalServerError)
        return
    }

    result := make(map[string]map[string]string)
    for _, ip := range ips {
        key := "device:" + ip
        vals, err := db.RDB.HGetAll(db.Ctx, key).Result()
        if err != nil {
            continue
        }
        result[ip] = vals
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
