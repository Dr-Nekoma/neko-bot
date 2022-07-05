package main

import (
	"fmt"
	"log"
	"neko-bot/BOT"
	"os"
	"os/signal"
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

func main() {

	token := os.Getenv("TOKEN")
	if token == "" {
		token = goDotEnvVariable("TOKEN")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discord.New("Bot " + token)

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

	messages := BOT.ParseCommand(m.Content)
	for i := 0; i < len(messages); i++ {
		s.ChannelMessageSendEmbed(m.ChannelID, jobMessage(messages[i]))
	}

}

func jobMessage(descr string) *discord.MessageEmbed {
	return &discord.MessageEmbed{
		URL:   "https://news.ycombinator.com/submitted?id=whoishiring",
		Title: "Job Hunting",
		Type:  "article",
		Color: 0x00acd7,
		Footer: &discord.MessageEmbedFooter{
			Text: "Made by Dr.Nekoma",
		},
		Image: &discord.MessageEmbedImage{
			URL: "https://jayclouse.com/wp-content/uploads/2019/06/hacker_news-1000x525-1.jpg",
		},
		Description: descr,
	}
}
