package main

import (
	"fmt"
	"log"
	"neko-bot/BOT"
	"neko-bot/DB"
	"neko-bot/MSG"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := env.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func init() {
	botToken = os.Getenv("TOKEN")
	if botToken == "" {
		botToken = goDotEnvVariable("TOKEN")
	}

	DB.ConnStr = os.Getenv("DATABASE_URL")
	if DB.ConnStr == "" {
		DB.ConnStr = goDotEnvVariable("URI")
	}

	port := os.Getenv("$PORT")
	if port != "" {
		strings.Replace(DB.ConnStr, goDotEnvVariable("PORT"), port, 1)
	}
}

var botToken string

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discord.New("Bot " + botToken)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discord.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func messageCreate(s *discord.Session, m *discord.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	messages := BOT.ParseCommand(m.Content, m.Author.Username)
	for i := 0; i < len(messages); i++ {
		message := messages[i]
		switch message.Kind {
		case MSG.Jobs:
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.JobMessage(messages[i].TitleLink, messages[i].Body))
		case MSG.LackOfJobs:
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.LackOfJobsMessage())
		case MSG.Project:
			switch message.SubKind {
			case MSG.ProjectAdd:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectAddMessage(messages[i], m.Author.Username))
			case MSG.ProjectList:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectListMessage(messages[i]))
			case MSG.ProjectEmptyList:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectEmptyListMessage())
			case MSG.ProjectDeleteId:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectDeleteIdMessage(messages[i], m.Author.Username))
			case MSG.ProjectDeleteIdea:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectDeleteIdeaMessage(messages[i], m.Author.Username))
			}
		case MSG.Error:
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.ErrorMessage(messages[i].Body))
		case MSG.Help:
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.HelpMessage(messages[i].Body))
		}
	}

}
