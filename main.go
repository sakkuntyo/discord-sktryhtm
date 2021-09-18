package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"log"
	"io"
	"time"
)

// 構造体定義
type Config struct {
	DiscordToken string `json:"discordToken"`
}

func main() {
	c := new(Config)
	jsonString, err := ioutil.ReadFile("settings.json")
	if err != nil {
		fmt.Println("error:\n", err)
		return
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		fmt.Println("error:\n", err)
		return
	}

	dg, err := discordgo.New("Bot " + c.DiscordToken)
	if err != nil {
		fmt.Println("error:start\n", err)
		return
	}
	dg.AddHandler(msgReceived)
	err = dg.Open()
	if err != nil {
		fmt.Println("error:wss\n", err)
		return
	}
	fmt.Println("BOT Running...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func msgReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	nickname := m.Author.Username
	member, err := s.State.Member(m.GuildID, m.Author.ID)

	if err == nil && member.Nick != "" {
		nickname = member.Nick
	}
	fmt.Println(m.Content + " by " + nickname)
	//fmt.Printf("(%%#v) %#v\n", m.Author)

	// Find a user's current voice channel
	vs, err := findUserVoiceState(s, m.Author.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// join the voice channel the user is on
	vc, err := s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// play audio file
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeFile("./test.mp3", opts)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	done := make(chan error)
	stream := dca.NewStream(encodeSession, vc, done)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}

			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}

	s.ChannelMessageSend(m.ChannelID, m.Content)
	fmt.Println("send message:", m.Content)
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}
