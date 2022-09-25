package handle

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error

type command struct {
	DiscordCommand *discordgo.ApplicationCommand
	Handler        commandHandler
}

func newCommand(name, description string, option []*discordgo.ApplicationCommandOption, handler commandHandler) *command {
	return &command{
		DiscordCommand: &discordgo.ApplicationCommand{
			ID:          name,
			Name:        name,
			Description: description,
		},
		Handler: handler,
	}
}

var commands = []*command{
	newCommand(
		"event",
		"Register an event at guild-level",
		[]*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "start-at",
				Description: `"now" or timestamp in format "yyyy-MM-dd hh:mm"`,
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "duration",
				Description: `Duration string, such as "15m", "1h", "1h30m". Defaults to "30m".`,
				Required:    false,
			},
		},
		func(s *discordgo.Session, i *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
			fmt.Printf("event:\n%+v\n\n", i)
			s.InteractionRespond(
				i,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "event!",
					},
				},
			)
			return nil
		},
	),
}
