package dmlock

import "encoding/json"

type Customer struct {
	ID                string `json:"id"`   // id
	LastHeartBeatNano int64  `json:"nano"` // last heart
}

func UnmarshalCustomerMap(s string) map[string]int64 {
	m := make(map[string]int64)
	_ = json.Unmarshal([]byte(s), &m)
	return m
}

func MarshalCustomerMap(m map[string]int64) string {
	raw, _ := json.Marshal(m)
	return string(raw)
}
