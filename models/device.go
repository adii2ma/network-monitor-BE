package models

type DeviceStatus struct {
    Name     string `json:"name"`
    IP       string `json:"ip"`
    Online   bool   `json:"online"`
    Location string `json:"location"`
    LastSeen int64  `json:"last_seen"`
    logs     string `json:"logs"` 
}
