package main
import (
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "github.com/bwmarrin/discordgo"
  "os/exec"
  "io/ioutil"
  "encoding/json"
)

// 構造体定義
type Config struct {
    DiscordToken  string  `json:"discordToken"`
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
  member,err := s.State.Member(m.GuildID, m.Author.ID)

  if err == nil && member.Nick != "" {
    nickname = member.Nick
  }
  fmt.Println(m.Content + " by " + nickname)

  if m.Content == "天気" || m.Content == "天気 全国" || m.Content == "天気　全国" {
    zenkokuTenkiNotify(s,m)
  }

  if m.Content == "天気 北海道" || m.Content == "天気　北海道" {
    hokkaidoTenkiNotify(s,m)
  }

  if m.Content == "天気 東北" || m.Content == "天気　東北" {
    tohokuTenkiNotify(s,m)
  }

  if m.Content == "天気 関東" || m.Content == "天気　関東" {
    kantoTenkiNotify(s,m)
  }

  if m.Content == "天気 北陸" || m.Content == "天気　北陸" {
    hokurikuTenkiNotify(s,m)
  }

  if m.Content == "天気 東海" || m.Content == "天気　東海" {
    tokaiTenkiNotify(s,m)
  }

  if m.Content == "天気 関西" || m.Content == "天気　関西" || m.Content == "天気 近畿" || m.Content == "天気　近畿" {
    kansaiTenkiNotify(s,m)
  }

  if m.Content == "天気 中国" || m.Content == "天気　中国" {
    chugokuTenkiNotify(s,m)
  }

  if m.Content == "天気 四国" || m.Content == "天気　四国" {
    shikokuTenkiNotify(s,m)
  }

  if m.Content == "天気 九州" || m.Content == "天気　九州" {
    kyusyuTenkiNotify(s,m)
  }

  if m.Content == "天気 沖縄" || m.Content == "天気　沖縄" || m.Content == "天気 那覇" || m.Content == "天気　那覇" {
    okinawaTenkiNotify(s,m)
  }
}

func zenkokuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"全国の天気予報")
  sendMsgRune,err := exec.Command("bash","./zenkoku-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func hokkaidoTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"北海道の天気予報")
  sendMsgRune,err := exec.Command("bash","./hokkaido-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func tohokuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"東北の天気予報")
  sendMsgRune,err := exec.Command("bash","./tohoku-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func kantoTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"関東の天気予報")
  sendMsgRune,err := exec.Command("bash","./kanto-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func hokurikuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"北陸の天気予報")
  sendMsgRune,err := exec.Command("bash","./hokuriku-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func tokaiTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"東海の天気予報")
  sendMsgRune,err := exec.Command("bash","./tokai-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func kansaiTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"関西の天気予報")
  sendMsgRune,err := exec.Command("bash","./kansai-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func chugokuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"中国の天気予報")
  sendMsgRune,err := exec.Command("bash","./chugoku-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func shikokuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"四国の天気予報")
  sendMsgRune,err := exec.Command("bash","./shikoku-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func kyusyuTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"九州の天気予報")
  sendMsgRune,err := exec.Command("bash","./kyusyu-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}

func okinawaTenkiNotify (s *discordgo.Session, m *discordgo.MessageCreate) {
  s.ChannelMessageSend(m.ChannelID,"沖縄の天気予報")
  sendMsgRune,err := exec.Command("bash","./okinawa-tenki.sh").Output()
  if err != nil {
    fmt.Println("error:", err)
  }
  sendMsg := string(sendMsgRune)
  fmt.Println(sendMsg)
  s.ChannelMessageSend(m.ChannelID,"```\n" + sendMsg + "\n```")
}
