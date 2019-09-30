package main

import (
	"encoding/json"
	"time"
)

var (
	lastMinute    int64 = 0
	lastTotal     int64 = 0
	lastHour      int64 = 0
	lastDay       int64 = 0
	currentMinute int64 = 0
	currentTotal  int64 = 0
	currentHour   int64 = 0
	currentDay    int64 = 0
)

var (
	mainTick = time.NewTicker(time.Second * 1)
)

func runMainTick() {
	for range mainTick.C {
		total, day, hour, minute := redisClient.getCounts()
		changed := false
		payload := &CountPayload{}
		if total != lastTotal {
			payload.Total = total
			lastTotal = total
			if !changed {
				changed = true
			}
		}
		if day != lastDay {
			payload.Day = day
			lastDay = day
			if !changed {
				changed = true
			}
		}
		if minute != lastMinute {
			payload.Minute = minute
			lastMinute = minute
			if !changed {
				changed = true
			}
		}
		if hour != lastHour {
			payload.Hour = hour
			lastHour = hour
			if !changed {
				changed = true
			}
		}
		currentMinute = minute
		currentDay = day
		currentHour = hour
		currentTotal = total
		if !changed {
			continue
		}
		bs, err := json.Marshal(payload)
		if err != nil {
			log.WithError(err).Error("Failed to marshal")
			return
		}
		hub.broadcast <- bs
	}
}
