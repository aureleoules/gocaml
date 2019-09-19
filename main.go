package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	d, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	d.AddHandler(onMessage)

	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	d.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	mention := "<@" + s.State.User.ID + ">"
	if strings.Contains(m.Content, mention) {

		reg := regexp.MustCompile("(?s)```(.*?)```")
		match := reg.FindStringSubmatch(m.Content)
		if match == nil {
			s.ChannelMessageSend(m.Message.ChannelID, "I don't understand...")
			return
		}
		code := match[1]
		result, err := evaluateCode(code)
		if err != nil {
			s.ChannelMessageSend(m.Message.ChannelID, "**ERROR**\n```"+err.Error()+"```")
			return
		}
		s.ChannelMessageSend(m.Message.ChannelID, "**Evaluation**:\n```"+strings.Replace(result, "        ", "", -1)+"```")
	}
}

func evaluateCode(code string) (string, error) {
	command := "echo \"" + code + "\" | ocaml"
	out, err := exec.Command("bash", "-c", command).Output()
	return string(out), err
}
