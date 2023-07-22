package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	token := loadToken()

	sess, err := discordgo.New(token)

	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		fmt.Printf("%+v", m.ChannelID)
		fmt.Println("")
		fmt.Printf("%+v", m.GuildID)
		fmt.Println("")

		if m.Content == "hi" {
			s.ChannelMessageSend(m.ChannelID, "hello")
		}

		//channel := m.ChannelID
		fmt.Printf("%+v", m.Author.ID)

		//s.ChannelVoiceJoin("1130105289031557222", "1130105289551646743", false, false)

		s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, false)

	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("The bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func loadToken() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading env variables %s", err)
	}

	token := os.Getenv("TOKEN")

	return token
}
