package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	Pipe  string
)

func init() {
	Token = os.Getenv("ALEX_TOKEN")
	Pipe = os.Getenv("ALEX_PIPE")
}

func main() {
	fmt.Println("Started...")
	fmt.Println(Token)
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents |= discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "$restart satisfactory" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Restarting satisfactory server...")
		fmt.Println("Command: $restart satisfactory")
		if err != nil {
			fmt.Println(err)
		}
		sendCmd("docker restart satisfactory-server")
	}
}

func sendCmd(cmd string) {
	file, err := os.OpenFile(Pipe, os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(cmd))
}
