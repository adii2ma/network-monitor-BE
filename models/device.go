package models

type DeviceStatus struct {
    IP       string `json:"ip"`
    Online   bool   `json:"online"`
    LastSeen int64  `json:"last_seen"`
}
