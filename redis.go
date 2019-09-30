package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	redisNS         = "pumpers"
	totalCacheCount = redisNS + "total"
)

type Redis struct {
	client *redis.Client
}

var redisClient = &Redis{}

func (r *Redis) connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Network:  "tcp",
		DB:       0,
		Password: "",
	})

	r.client = client

	err := client.Ping().Err()
	if err != nil {
		panic(err)
	}

	log.Info("Connected to Redis.")
}

func (r *Redis) getCounts() (int64, int64, int64, int64) {
	now := time.Now()
	dayTimestamp := now.Format("2006-01-01")
	hour := now.Hour()
	min := now.Minute()

	totalKey := fmt.Sprintf("%s::total", redisNS)
	dayKey := fmt.Sprintf("%s::day::%s", redisNS, dayTimestamp)
	hourKey := fmt.Sprintf("%s::day-hour::%s::hr::%d", redisNS, dayTimestamp, hour)
	minKey := fmt.Sprintf("%s::day-hour-minute::%s::hr::%d::min::%d", redisNS, dayTimestamp, hour, min)

	total, _ := r.client.Get(totalKey).Int64()
	dy, err := r.client.Get(dayKey).Int64()
	if err != nil {
		log.WithError(err)
	}
	hr, err := r.client.Get(hourKey).Int64()
	if err != nil {
		log.WithError(err)
	}
	mn, err := r.client.Get(minKey).Int64()

	if err != nil {
		log.WithError(err)
	}

	return total, dy, hr, mn
}

func (r *Redis) increment(amount int) (int64, int64, int64, int64) {
	amount64 := int64(amount)

	now := time.Now()
	dayTimestamp := now.Format("2006-01-01")
	hour := now.Hour()
	min := now.Minute()

	totalKey := fmt.Sprintf("%s::total", redisNS)
	dayKey := fmt.Sprintf("%s::day::%s", redisNS, dayTimestamp)
	hourKey := fmt.Sprintf("%s::day-hour::%s::hr::%d", redisNS, dayTimestamp, hour)
	minKey := fmt.Sprintf("%s::day-hour-minute::%s::hr::%d::min::%d", redisNS, dayTimestamp, hour, min)

	total := r.client.IncrBy(totalKey, amount64).Val()
	dy := r.client.IncrBy(dayKey, amount64).Val()
	hr := r.client.IncrBy(hourKey, amount64).Val()
	mn := r.client.IncrBy(minKey, amount64).Val()

	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	endOfNextMonth := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 31)

	// Expire minute key stamp at end of day
	_ = r.client.ExpireAt(minKey, endOfDay).Err()
	_ = r.client.ExpireAt(hourKey, endOfNextMonth)
	_ = r.client.ExpireAt(dayKey, endOfNextMonth).Err()

	return total, dy, hr, mn
}
