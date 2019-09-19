package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"time"

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

		reg := regexp.MustCompile("(?s)```(ocaml)?(.*?)```")
		match := reg.FindStringSubmatch(m.Content)
		if match == nil {
			s.ChannelMessageSend(m.Message.ChannelID, "I don't understand...")
			return
		}
		code := match[len(match)-1]
		result, err := evaluateCode(code)
		if err != nil {
			s.ChannelMessageSend(m.Message.ChannelID, "**ERROR**\n```"+err.Error()+"```")
			return
		}

		//format
		formatted := strings.Replace(result, "        ", "", -1)
		formatted = strings.Replace(formatted, "   ", " ", -1)
		formatted = removeLastLine(formatted)
		s.ChannelMessageSend(m.Message.ChannelID, "**Evaluation**:\n```ocaml\n"+formatted+"```")
	}
}

func evaluateCode(code string) (string, error) {
	command := "echo \"" + code + "\" | ocaml"
	process := exec.Command("bash", "-c", command)

	go func() {
		time.Sleep(3 * time.Second)
		terminateProc(process.Process.Pid, os.Interrupt)
	}()

	out, err := process.Output()
	return string(out), err
}
