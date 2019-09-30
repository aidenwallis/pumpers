package main

type Counter struct{}

type CountPayload struct {
	Total  int64 `json:"t,omitempty"`
	Day    int64 `json:"d,omitempty"`
	Hour   int64 `json:"h,omitempty"`
	Minute int64 `json:"m,omitempty"`
}

var counter = &Counter{}

func (c *Counter) increment(amount int) {
	total, day, hour, minute := redisClient.increment(amount)
	currentTotal = total
	currentDay = day
	currentHour = hour
	currentMinute = minute
}
