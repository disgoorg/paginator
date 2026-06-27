package main

import (
	"context"
	"log/slog"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/paginator"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "test",
			Description: "simple test command",
		},
	}
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.Info("starting example...")
	slog.Info("disgo version: " + disgo.Version)

	manager := paginator.New()
	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(commandListener(manager)),
		bot.WithEventListeners(manager),
	)
	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest.SetGuildCommands(client.ApplicationID, guildID, commands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(manager *paginator.Manager) func(event *events.ApplicationCommandInteractionCreate) {
	return func(event *events.ApplicationCommandInteractionCreate) {
		data := event.SlashCommandInteractionData()
		if data.CommandName() == "test" {
			pData := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
			err := manager.Create(event.Respond, paginator.Pages{
				ID: event.ID().String(),
				PageFunc: func(page int, embed discord.Embed) discord.Embed {
					description := ""
					for i := 0; i < 5; i++ {
						if page*5+i >= len(pData) {
							break
						}
						description += pData[page*5+i] + "\n"
					}
					return embed.WithTitlef("Page %d", page).WithDescription(description)
				},
				Pages:      int(math.Ceil(float64(len(pData)) / 5)),
				Creator:    event.User().ID,
				ExpireMode: paginator.ExpireModeAfterLastUsage,
			}, false)
			if err != nil {
				slog.Error("command error", slog.Any("err", err))
			}
		}
	}
}
