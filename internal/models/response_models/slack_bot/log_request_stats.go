package slack_bot

import "time"

type LogRequestStats struct {
	RequestCount int       `json:"requestCount"`
	Time         time.Time `json:"time"`
}
