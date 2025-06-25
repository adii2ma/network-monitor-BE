package models

type DeviceStatus struct {
    IP       string `json:"ip"`
    Online   bool   `json:"online"`
    Location string `json:"location"`
    LastSeen int64  `json:"last_seen"`
}
