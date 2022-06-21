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

	token := goDotEnvVariable("TOKEN")

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

	// String manipulation: Slice string

	// neko!jobs 3 Haskell

	// If the message is "ping" reply with "Pong!"
	if m.Content == "neko!jobs 3 Haskell" {
		messages := BOT.HackerNewsJobs("Haskell", 3)
		for i := 0; i < 3; i++ {
			s.ChannelMessageSend(m.ChannelID, messages[i])
		}
	}

}
