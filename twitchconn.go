package main

import (
	twitch "github.com/gempir/go-twitch-irc/v2"
)

const emoteID = "303636564"

func startConnection() {
	// connect to twitch anonymously
	client := twitch.NewClient("justinfan123", "oauth:123123123123")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, emote := range message.Emotes {
			if emote.ID == emoteID {
				counter.increment(emote.Count)
				return
			}
		}
	})

	client.Join("xqcow")

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
