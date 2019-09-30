package main

import "encoding/json"

type Counter struct{}

type CountPayload struct {
	Total  int64 `json:"t"`
	Day    int64 `json:"d"`
	Hour   int64 `json:"h"`
	Minute int64 `json:"m"`
}

var counter = &Counter{}

func (c *Counter) increment(amount int) {
	total, day, hour, minute := redisClient.increment(amount)
	payload := &CountPayload{
		Total:  total,
		Day:    day,
		Hour:   hour,
		Minute: minute,
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Failed to marshal")
		return
	}
	hub.broadcast <- bs
}
