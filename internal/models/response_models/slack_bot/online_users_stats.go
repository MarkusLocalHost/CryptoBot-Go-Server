package slack_bot

import "time"

type OnlineUserStats struct {
	OnlineUsers int       `json:"online_users"`
	IdsUsers    []int64   `json:"ids_users"`
	Time        time.Time `json:"time"`
}
