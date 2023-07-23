package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Andreychik32/ytdl"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	token := loadToken()

	sess, err := discordgo.New(token)

	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(onMessageCreate)

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

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate)  {
		if m.Author.ID == s.State.User.ID || m.Author.Bot {
			return
		}

		// Check if the message starts with the command prefix, such as "!play".
		if !strings.HasPrefix(m.Content, "!play") {
			return
		}

		// Extract the YouTube URL from the message content.
		youtubeURL := strings.TrimSpace(strings.TrimPrefix(m.Content, "!play"))

		// Play the song from the YouTube URL.
		err := playFromYouTubeURL(s, m.GuildID, m.Author.ID, youtubeURL)
		if err != nil {
			fmt.Println("Error playing song: ", err)
			return
		}

		if m.Content == "hi" {
			s.ChannelMessageSend(m.ChannelID, "hello")

			connectToVoiceChannel(s, m.GuildID, m.Author.ID)
	}
}

func connectToVoiceChannel(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceConnection, error) {
	// Find the user in the guild's voice state.
	vs, err := s.State.VoiceState(guildID, userID)
	if err != nil {
		return nil, err
	}

	// Connect to the voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, vs.ChannelID, false, false)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func playFromYouTubeURL(s *discordgo.Session, guildID, userID, url string) error {
	// Connect to the user's voice channel.
	vc, err := connectToVoiceChannel(s, guildID, userID)
	if err != nil {
		return err
	}
	defer vc.Disconnect()

	// Download the YouTube audio from the given URL.
	audioFile, err := downloadYouTubeAudio(url)
	if err != nil {
		return err
	}
	defer os.Remove(audioFile)

	// Start playing the audio in the voice channel.
	err = playAudio(vc, audioFile)
	if err != nil {
		return err
	}

	return nil
}

func downloadYouTubeAudio(url string) (string, error) {
	// Create a new video info object.
	ctx := context.Background()
	client := ytdl.DefaultClient
	videoInfo, err := ytdl.GetVideoInfo(ctx, "https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	if err != nil {
		return "", err
	}

	// Download the audio stream.
	audioFile, err := os.Create(videoInfo.Title + ".mp4")
	if err != nil {
		return "", err
	}
	defer audioFile.Close()

	err = client.Download(ctx, videoInfo, videoInfo.Formats[0], audioFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", err)

	return videoInfo.Title + ".mp4", err
}

func playAudio(vc *discordgo.VoiceConnection, audioFile string) error {
	// Open the audio file.
	file, err := os.Open(audioFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start speaking.
	vc.Speaking(true)

	// Send the audio packets.
	dgvoice.PlayAudioFile(vc, "test.mp3", make(chan bool))

	// Stop speaking.
	vc.Speaking(false)

	return nil
}
