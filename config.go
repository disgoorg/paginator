package paginator

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

func DefaultConfig() *Config {
	return &Config{
		ButtonsConfig: ButtonsConfig{
			First: &ComponentOptions{
				Emoji: discord.ComponentEmoji{
					Name: "⏮",
				},
				Style: discord.ButtonStylePrimary,
			},
			Back: &ComponentOptions{
				Emoji: discord.ComponentEmoji{
					Name: "◀",
				},
				Style: discord.ButtonStylePrimary,
			},
			Stop: &ComponentOptions{
				Emoji: discord.ComponentEmoji{
					Name: "🗑",
				},
				Style: discord.ButtonStyleDanger,
			},
			Next: &ComponentOptions{
				Emoji: discord.ComponentEmoji{
					Name: "▶",
				},
				Style: discord.ButtonStylePrimary,
			},
			Last: &ComponentOptions{
				Emoji: discord.ComponentEmoji{
					Name: "⏩",
				},
				Style: discord.ButtonStylePrimary,
			},
		},
		NoPermissionMessage: "You can't interact with this paginator because it's not yours.",
		CustomIDPrefix:      "paginator",
		EmbedColor:          0x4c50c1,
		CleanupInterval:     30 * time.Second,
		ExpireTime:          5 * time.Minute,
	}
}

type Config struct {
	ButtonsConfig       ButtonsConfig
	NoPermissionMessage string
	CustomIDPrefix      string
	EmbedColor          int
	CleanupInterval     time.Duration
	ExpireTime          time.Duration
}

type ButtonsConfig struct {
	First *ComponentOptions
	Back  *ComponentOptions
	Stop  *ComponentOptions
	Next  *ComponentOptions
	Last  *ComponentOptions
}

type ComponentOptions struct {
	Emoji discord.ComponentEmoji
	Label string
	Style discord.ButtonStyle
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithButtonsConfig(buttonsConfig ButtonsConfig) ConfigOpt {
	return func(config *Config) {
		config.ButtonsConfig = buttonsConfig
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithNoPermissionMessage(noPermissionMessage string) ConfigOpt {
	return func(config *Config) {
		config.NoPermissionMessage = noPermissionMessage
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCustomIDPrefix(prefix string) ConfigOpt {
	return func(config *Config) {
		config.CustomIDPrefix = prefix
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithEmbedColor(color int) ConfigOpt {
	return func(config *Config) {
		config.EmbedColor = color
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithCleanupInterval(cleanupInterval time.Duration) ConfigOpt {
	return func(config *Config) {
		config.CleanupInterval = cleanupInterval
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithTimeout(timeout time.Duration) ConfigOpt {
	return func(config *Config) {
		config.ExpireTime = timeout
	}
}
