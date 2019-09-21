package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"syscall"

	"github.com/aureleoules/gocaml/db"
	"github.com/aureleoules/gocaml/models"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

var prefix = "!gocaml"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	db.Connect(os.Getenv("URI"), os.Getenv("DATABASE"))

	d, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	d.AddHandler(onMessage)
	d.AddHandler(onMessageUpdate)

	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	d.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// don't process bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	isStats, userID := IsStats(m)
	if isStats {
		if userID == "" {
			users, err := models.GetUsers()
			if err != nil {
				log.Println(err)
			}

			s.ChannelMessageSend(m.Message.ChannelID, ParseStats(users))
			return
		}
		user, err := models.GetUser(userID)
		if err != nil {
			log.Println(err)
		}
		s.ChannelMessageSend(m.Message.ChannelID, ParseStats([]models.User{user}))

	}

	isEval, code := IsCodeEvaluation(m)
	if !isEval {
		return
	}

	result, err := evaluateCode(code)
	if err != nil {
		s.ChannelMessageSend(m.Message.ChannelID, "**RUNTIME ERROR**\n```"+err.Error()+"```")
		return
	}

	//format
	var formatted string
	if result != "" {
		formatted = FormatEvaluation(result)
	}

	s.ChannelMessageSend(m.Message.ChannelID, "**Evaluation**:\n```ocaml\n"+formatted+"```")

	user, err := models.GetUser(m.Author.ID)
	if err == mgo.ErrNotFound {
		u := models.User{
			DiscordID:     m.Author.ID,
			Username:      m.Author.Username,
			Discriminator: m.Author.Discriminator,
		}
		user, err = u.Create()
		if err != nil {
			log.Println(err)
		}
	}
	if ContainsError(formatted) {
		user.IncrementError()
	} else {
		user.IncrementSuccess()
	}
}

func onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	message := discordgo.MessageCreate(*m)
	onMessage(s, &message)
}

func evaluateCode(code string) (string, error) {
	command := "echo \"" + code + "\" | ocaml"
	process := exec.Command("bash", "-c", command)

	terminated := false
	go func() {
		time.Sleep(5 * time.Second)
		if terminated {
			return
		}
		p := exec.Command("bash", "-c", "pkill -f ocamlrun")
		_, err := p.Output()
		if err != nil {
			log.Println(err)
		}
	}()
	out, err := process.Output()
	terminated = true
	return string(out), err
}
