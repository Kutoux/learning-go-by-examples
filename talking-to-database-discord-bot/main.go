package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Answers struct {
	//keeps track of original channel for response
	OriginChannelId string
	FavFood         string
	FavGame         string
	RecordID        int64
}

func (a *Answers) ToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Favorite Food",
			Value: a.FavFood,
		},
		{
			Name:  "Favorite Game",
			Value: a.FavGame,
		},
		{
			Name:  "Record Id",
			Value: strconv.FormatInt(a.RecordID, 10), //converts int64 into string
		},
	}

	return discordgo.MessageEmbed{
		Title:  "New responses!",
		Fields: fields,
	}
}

var responses map[string]Answers = map[string]Answers{}

// subcommands for prefix
const prefix string = "!gobot"

func main() {
	godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")

	sess, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Fatal(err)
	}

	//Connecting to database
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		//DM Logic
		if m.GuildID == "" { //If so, then it's a DM
			userPromptResponseHandler(db, s, m)
		}

		//Server Logic
		args := strings.Split(m.Content, " ")

		//if message doesn't start with prefix command, ignore
		if args[0] != prefix {
			return
		}

		if args[1] == "hello" {
			HelloWorldHandler(s, m)
		}

		if args[1] == "vow" {
			DwarvenVowsHandler(s, m)
		}

		if args[1] == "prompt" {
			UserPromptHandler(s, m)
		}

		if args[1] == "answers" {
			AnswersHandler(db, s, m)
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

func AnswersHandler(db *sql.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	spl := strings.Split(m.Content, " ")
	if len(spl) > 3 {
		s.ChannelMessageSend(m.ChannelID, "An ID must be provided. Ex: '!gobot answers 1'")
		return
	}

	//convert 3rd element into an int
	id, err := strconv.Atoi(spl[2])
	if err != nil {
		log.Fatal(err)
	}

	var recordId int64
	var answerStr string
	var userId int64

	query := "select * from discord_messages where id = ?"
	row := db.QueryRow(query, id)
	err = row.Scan(&recordId, &answerStr, &userId) //taking values from the row in that order and placing into variables
	if err != nil {
		log.Fatal(err)
	}

	var answers Answers                               //converts answer string into Answers object
	err = json.Unmarshal([]byte(answerStr), &answers) //converts string back into byte slice/byte array, and pass &answers so it knows where to put it
	if err != nil {
		log.Fatal(err)
	}

	answers.RecordID = recordId //update answers object with recordId
	embed := answers.ToMessageEmbed()
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func userPromptResponseHandler(db *sql.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	answers, ok := responses[m.ChannelID] //should correspond within our map
	if !ok {
		//if someone DMs bot and the bot isn't tracking a question, just ignore
		return
	}

	if answers.FavFood == "" { //If empty, we're still waiting on this question
		answers.FavFood = m.Content

		s.ChannelMessageSend(m.ChannelID, "Great! Now, what's your favorite game?")

		responses[m.ChannelID] = answers //updates response

		return
	} else {
		//if FavFood has content, move onto next question
		answers.FavGame = m.Content
		//log.Printf("answers: %v, %v", answers.FavFood, answers.FavGame)

		query := "insert into discord_messages (payload, user_id) values (?,?)"
		jbytes, err := json.Marshal(answers) //conversion of answers into byte array
		if err != nil {
			log.Fatal(err)
		}

		//executes query, returns two objects
		res, err := db.Exec(query, string(jbytes), m.ChannelID) //parameters(query, ?(payload column), ?(unique identifier)) from above
		if err != nil {
			log.Fatal(err)
		}

		//get id of last column created
		lastInserted, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		answers.RecordID = lastInserted

		embed := answers.ToMessageEmbed()
		s.ChannelMessageSendEmbed(answers.OriginChannelId, &embed)

		//closes responses
		delete(responses, m.ChannelID)
	}
}

func DwarvenVowsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	author := discordgo.MessageEmbedAuthor{
		Name: "Lloyd Irving",
	}

	embed := discordgo.MessageEmbed{
		Title:  dwarvenVows[selection],
		Author: &author,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func HelloWorldHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "world!")
}

func UserPromptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//User Channel: DM between bot and user
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Panic(err)
	}
	//If the user is already answering questions, ignore it, otherwise ask questions
	if _, ok := responses[channel.ID]; !ok {
		responses[channel.ID] = Answers{
			OriginChannelId: m.ChannelID,
			FavFood:         "",
			FavGame:         "",
		}
		s.ChannelMessageSend(channel.ID, "Hello! Here are a few questions:")
		s.ChannelMessageSend(channel.ID, "What's your favorite food?")
	} else {
		s.ChannelMessageSend(channel.ID, "I'm still waiting on your response. :(")
	}
}
