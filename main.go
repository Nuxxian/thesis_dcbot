package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token     = os.Getenv("TOKEN")
	channelID = os.Getenv("CHANNEL_ID")
)

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Schedule the daily message
	go func() {
		for {
			postDailyMessage(dg)
			time.Sleep(24 * time.Hour)
		}
	}()

	// Wait here until CTRL+C or other term signal is received
	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func postDailyMessage(s *discordgo.Session) {
	deadline := time.Date(2024, time.May, 23, 00, 00, 00, 0, time.Now().Location())
	timeLeft := int(deadline.Sub(time.Now()).Hours() / 24)
	data := strconv.Itoa(timeLeft)
	_, err := s.ChannelMessageSend(channelID, "Nog "+data+" dagen!")
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}
