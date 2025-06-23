package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "backend/store"
    "backend/db"
    "time"
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
        if err != nil || len(vals) == 0 {
            // Provide default status if missing
            vals = map[string]string{
                "online":    "false",
                "last_seen": "0",
            }
        }
        result[ip] = vals
        // Log status changes
        logStatusChange(ip, vals)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func LogsHandler(w http.ResponseWriter, r *http.Request) {
    logs, err := db.RDB.LRange(db.Ctx, "logs", 0, 99).Result() // Get last 100 logs
    if err != nil {
        http.Error(w, "Failed to get logs", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "logs": logs,
    })
}

func logStatusChange(ip string, status map[string]string) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    online := status["online"]
    lastSeen := status["last_seen"]
    
    logEntry := fmt.Sprintf("[%s] %s: %s (Last seen: %s)", 
        timestamp, ip, 
        map[string]string{"true": "Online", "false": "Offline"}[online],
        time.Unix(parseInt64(lastSeen), 0).Format("15:04:05"))
    
    // Add to logs list (keep only last 1000 entries)
    db.RDB.LPush(db.Ctx, "logs", logEntry)
    db.RDB.LTrim(db.Ctx, "logs", 0, 999)
}

func parseInt64(s string) int64 {
    var i int64
    fmt.Sscanf(s, "%d", &i)
    return i
}
