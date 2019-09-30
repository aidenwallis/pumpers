package main

import "github.com/sirupsen/logrus"

var log = logrus.StandardLogger()

func main() {
	redisClient.connect()
	go hub.run()
	go startConnection()

	log.Info("Starting ResidentSleeper...")
	startServer()
}
