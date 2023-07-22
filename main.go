package main

import (
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
		guildID := m.GuildID

		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "hi" {
			s.ChannelMessageSend(m.ChannelID, "hello")
			voiceState, err := s.State.VoiceState(guildID, m.Author.ID)

			if err != nil {
				log.Fatal(err)
			}

			authorChannelID := voiceState.ChannelID
			s.ChannelVoiceJoin(guildID, authorChannelID, false, false)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAll

	err = sess.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

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
