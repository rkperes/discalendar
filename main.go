package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	AppID string
)

func init() {
	flag.StringVar(&AppID, "a", "", "Application ID")
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	_, err = dg.ApplicationCommandCreate(AppID, "", &discordgo.ApplicationCommand{
		ID:            "ping",
		Name:          "ping",
		Type:          discordgo.ChatApplicationCommand,
		ApplicationID: AppID,
		Description:   "ping",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(applicationCommand)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func applicationCommand(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	interaction := ic.Interaction
	if interaction.Type == discordgo.InteractionApplicationCommand {
		data := interaction.ApplicationCommandData()
		switch data.Name {
		case "ping":
			s.InteractionRespond(
				interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "pong!",
					},
				},
			)
		default:
			log.Printf("invalid command %q", data.Name)
		}
	}
}
