package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main()  {
	sess, err := discordgo.New("Bot MTA5OTcyOTkxNjQwNjI4ODQ5NQ.GsNQkL.lM1T7l84U2Q7WJ8CjEFLAe1ea2JOa8rMdRN9HQ")

	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func (s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "hi" {
			s.ChannelMessageSend(m.ChannelID, "hello")
		}
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