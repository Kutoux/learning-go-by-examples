package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// subcommands for prefix
const prefix string = "!gobot"

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	sess, err := discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		args := strings.Split(m.Content, " ")

		//if message doesn't start with prefix command, ignore
		if args[0] != prefix {
			return
		}

		if args[1] == "hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}

		if args[1] == "dwarven" && args[2] == "vow" {
			dwarvenVows := []string{
				"Dwarven Vow #1: Let's all work together for a peaceful world.",
				"Dwarven Vow #2: Never abandon someone in need.",
				"Dwarven Vow #4: Don't depend on others. Walk on your own two legs.",
				"Dwarven Vow #5: The more you add, the worse it gets.",
				"Dwarven Vow #7: Justice and love will always win.",
				"Dwarven Vow #7: Goodness and love will always win.",
				"Dwarven Vow #9: Fall down seven times, stand up eight.",
				"Dwarven Vow #10: Play hard, play often.",
				"Dwarven Vow #11: Lying is the first step down the path of thievery.",
				"Dwarven Vow #14: Even a small star shines in the darkness.",
				"Dwarven Vow #16: You can do anything if you try.",
				"Dwarven Vow #18: It's better to be deceived than to deceive.",
				"Dwarven Vow #24: Never let your feet run faster than your shoes.",
				"Dwarven Vow #32: Cross even a stone bridge after you've tested it.",
				"Dwarven Vow #41: It's better to begin in the evening than not at all.",
				"Dwarven Vow #41: Haste makes waste.",
				"Dwarven Vow #43: Never forget the basics.",
				"Dwarven Vow #55: A bad workman blames his tools.",
				"Dwarven Vow #108: Let sleeping dogs lie.",
				"Dwarven Vow #134: Compassion benefits all men.",
			}

			selection := rand.Intn(len(dwarvenVows))

			// author := discordgo.MessageEmbedAuthor{
			// 	Name: "Dirk",
			// }

			embed := discordgo.MessageEmbed{
				Title: dwarvenVows[selection],
				//Author: &author,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("The bot is online.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
