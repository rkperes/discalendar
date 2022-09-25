package handle

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO: find a better name
type Controller struct {
	dg       *discordgo.Session
	appID    string
	token    string
	commands map[string]*command
}

func NewController(appID string) *Controller {
	commandsByName := make(map[string]*command, len(commands))
	for _, command := range commands {
		commandsByName[command.DiscordCommand.Name] = command
		command.DiscordCommand.ApplicationID = appID
	}

	return &Controller{
		appID:    appID,
		commands: commandsByName,
	}
}

func (c *Controller) RegisterCommands(dg *discordgo.Session) error {
	c.dg = dg
	_ = registerCommands(c.dg, c.appID, c.commands)
	dg.AddHandler(handler(c.commands))
	return nil
}

func registerCommands(s *discordgo.Session, appID string, commands map[string]*command) error {
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(appID, "", command.DiscordCommand)
		if err != nil {
			return fmt.Errorf("could not register command %q: %w", command.DiscordCommand.Name, err)
		}
	}
	return nil
}

func handler(commands map[string]*command) func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, ic *discordgo.InteractionCreate) {
		interaction := ic.Interaction
		if interaction.Type != discordgo.InteractionApplicationCommand {
			log.Printf("invalid interaction %q", interaction.Type)
			return
		}
		data := interaction.ApplicationCommandData()
		command, ok := commands[data.Name]
		if !ok {
			log.Printf("invalid command %q", data.Name)
			return
		}
		command.Handler(s, interaction, data)
	}
}
