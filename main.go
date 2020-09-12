package main
import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "github.com/bwmarrin/discordgo"
  "os/exec"
)

const(
  TOKEN = ""
)

func main() {
  dg, err := discordgo.New("Bot " + TOKEN)
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
  member,err := s.State.Member(m.GuildID, m.Author.ID)

  if err == nil && member.Nick != "" {
    nickname = member.Nick
  }
  fmt.Println(m.Content + " by " + nickname)

  if m.Content == "hi" {
    s.ChannelMessageSend(m.ChannelID, "Hi!" + nickname)
  }

  if m.Content == "天気 関東" {
    //関東
    s.ChannelMessageSend(m.ChannelID,"関東の天気予報")
    sendMsgRune,err := exec.Command("bash","/root/otenkimaru/kanto-tenki.sh").Output()
    if err != nil {
      fmt.Println("error:exec failed")
    }
    sendMsg := string(sendMsgRune)
    fmt.Println(sendMsg)
    s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
  }

  if m.Content == "天気 関西" {
    //関西
    s.ChannelMessageSend(m.ChannelID,"関西の天気予報")
    sendMsgRune,err := exec.Command("bash","/root/otenkimaru/kansai-tenki.sh").Output()
    if err != nil {
      fmt.Println("error:exec failed")
    }
    sendMsg := string(sendMsgRune)
    fmt.Println(sendMsg)
    s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
  }
}
