package main

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"github.com/disgoorg/paginator"
	"github.com/disgoorg/snowflake/v2"
	"math"
	"os"
	"os/signal"
	"syscall"
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
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	manager := paginator.New()
	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(commandListener(manager)),
		bot.WithEventListeners(manager),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
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
				PageFunc: func(page int, embed *discord.EmbedBuilder) {
					embed.SetTitlef("Page %d", page+1)

					description := ""
					for i := 0; i < 5; i++ {
						if page*5+i >= len(pData) {
							break
						}
						description += pData[page*5+i] + "\n"
					}
					embed.SetDescription(description)
				},
				Pages:      int(math.Ceil(float64(len(pData)) / 5)),
				Creator:    event.User().ID,
				ExpireMode: paginator.ExpireModeAfterLastUsage,
			}, false)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
