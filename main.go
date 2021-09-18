package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
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
	// It works only with the rythm channel of the dust box.
	if m.ChannelID != "888667445999587328" {
		return
	}

	// start play
	if regexp.MustCompile(`^!sp `).Match([]byte(m.Content)) {
		musicInfo := strings.Replace(m.Content, "!sp ", "", 1)
		s.ChannelMessageSend(m.ChannelID, "再生を試みるんゴ")

		nickname := m.Author.Username
		member, err := s.State.Member(m.GuildID, m.Author.ID)

		if err == nil && member.Nick != "" {
			nickname = member.Nick
		}
		fmt.Println(m.Content + " by " + nickname)
		//fmt.Printf("(%%#v) %#v\n", m.ChannelID)

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

		// dl youtube audio
		cmdstr := "youtube-dl 'ytsearch:"
		cmdstr += musicInfo // moviename
		cmdstr += "' -x --audio-format mp3 -o './tmp/%(title)s.%(ext)s' |  grep -e 'mp3' -e ': Downloading webpage' | sed -e \"s/\\[youtube\\] /https:\\/\\/www\\.youtube\\.com\\/watch\\?v=/g\" -e \"s/: Downloading webpage//g\" -e \"s/.*Destination: //g\""

		out, err := exec.Command("sh", "-c", cmdstr).Output()
		if err != nil {
			println("exec.Commandに失敗")
			s.ChannelMessageSend(m.ChannelID, "曲のDLに失敗したんゴ")
			return
		}

		splitOut := regexp.MustCompile(`\r\n|\n`).Split(string(out), -1)
		movieUrl := splitOut[0]
		audioFilePath := splitOut[1]

		// play audio file
		opts := dca.StdEncodeOptions
		opts.RawOutput = true
		opts.Bitrate = 120
		opts.Volume = 1

		encodeSession, err := dca.EncodeFile(audioFilePath, opts)
		if err != nil {
			log.Fatal("Failed creating an encoding session: ", err)
			s.ChannelMessageSend(m.ChannelID, "エンコードに失敗したんゴ")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "play -> "+movieUrl)

		done := make(chan error)
		stream := dca.NewStream(encodeSession, vc, done)

		ticker := time.NewTicker(time.Second)

		for {
			select {
			case err := <-done:
				if err != nil && err != io.EOF {
					log.Fatal("An error occured", err)
					s.ChannelMessageSend(m.ChannelID, "よくわからん処理に失敗したんゴ、すまんゴ")
					return
				}

				// Clean up incase something happened and ffmpeg is still running
				encodeSession.Truncate()
				s.ChannelMessageSend(m.ChannelID, "再生したから落ちるんゴ、と思ったけど抜け方わからないんゴ")
				return
			case <-ticker.C:
				stats := encodeSession.Stats()
				playbackPosition := stream.PlaybackPosition()

				fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
			}
		}
	}
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
