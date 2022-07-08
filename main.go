package main

import (
	"fmt"
	"log"
	"neko-bot/BOT"
	"neko-bot/DB"
	"neko-bot/MSG"
	"os"
	"os/signal"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
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

}

var botToken string

func main() {

	dg, err := discord.New("Bot " + botToken)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discord.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	err = dg.UpdateListeningStatus("neko!help")
	if err != nil {
		fmt.Println("error with status,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

}

func messageCreate(s *discord.Session, m *discord.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	messages := BOT.ParseCommand(m.Content, m.Author.Username)
	if len(messages) == 0 {
		return
	}
	if messages[0].Kind == MSG.Project && messages[0].SubKind == MSG.ProjectList {
		msg, _ := s.ChannelMessageSend(m.ChannelID, ("Hey " + m.Author.Mention() + "!"))
		ch, err := s.MessageThreadStart(m.ChannelID, msg.ID, MSG.ProjectListTitle, 60)
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.ErrorMessage("I couldn't create a thread for projects!"))
			return
		}
		for i := 0; i < len(messages); i++ {
			s.ChannelMessageSendEmbed(ch.ID, MSG.ProjectListMessage(messages[i]))
		}
		return
	}
	if messages[0].Kind == MSG.Jobs {
		msg, _ := s.ChannelMessageSend(m.ChannelID, ("Hey " + m.Author.Mention() + "!"))
		ch, err := s.MessageThreadStart(m.ChannelID, msg.ID, MSG.JobsTitle, 1440)
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.ErrorMessage("I couldn't create a thread for jobs!"))
			return
		}
		for i := 0; i < len(messages); i++ {
			s.ChannelMessageSendEmbed(ch.ID, MSG.JobMessage(messages[i].TitleLink, messages[i].Body))
		}
		return
	}

	for i := 0; i < len(messages); i++ {
		message := messages[i]
		switch message.Kind {
		case MSG.LackOfJobs:
			s.ChannelMessageSendEmbed(m.ChannelID, MSG.LackOfJobsMessage())
		case MSG.Project:
			switch message.SubKind {
			case MSG.ProjectAdd:
				s.ChannelMessageSendEmbed(m.ChannelID, MSG.ProjectAddMessage(messages[i], m.Author.Username))
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
