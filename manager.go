package paginator

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var _ bot.EventListener = (*Manager)(nil)

type ExpireMode int

const (
	ExpireModeAfterCreation ExpireMode = iota
	ExpireModeAfterLastUsage
)

type Pages struct {
	ID          string
	PageFunc    func(page int, embed *discord.EmbedBuilder)
	Pages       int
	Creator     snowflake.ID
	ExpireMode  ExpireMode
	lastUsed    time.Time
	currentPage int
}

func New(opts ...ConfigOpt) *Manager {
	config := DefaultConfig()
	config.Apply(opts)
	paginator := &Manager{
		config: *config,
		pages:  map[string]*Pages{},
	}
	go paginator.startCleanup()
	return paginator
}

type Manager struct {
	config  Config
	pagesMu sync.Mutex
	pages   map[string]*Pages
}

func (m *Manager) startCleanup() {
	ticker := time.NewTicker(m.config.CleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		m.cleanup()
	}
}

func (m *Manager) cleanup() {
	m.pagesMu.Lock()
	defer m.pagesMu.Unlock()

	now := time.Now()
	for id, paginator := range m.pages {
		if now.Sub(paginator.lastUsed) > m.config.ExpireTime {
			delete(m.pages, id)
		}
	}
}

func (m *Manager) Update(responderFunc events.InteractionResponderFunc, pages Pages) error {
	m.add(&pages)

	return responderFunc(discord.InteractionResponseTypeUpdateMessage, m.makeMessageUpdate(&pages))
}

func (m *Manager) Create(responderFunc events.InteractionResponderFunc, pages Pages, ephemeral bool) error {
	m.add(&pages)

	return responderFunc(discord.InteractionResponseTypeCreateMessage, m.makeMessageCreate(&pages, ephemeral))
}

func (m *Manager) UpdateMessage(client bot.Client, channelID snowflake.ID, messageID snowflake.ID, pages Pages) (*discord.Message, error) {
	m.add(&pages)

	return client.Rest().UpdateMessage(channelID, messageID, m.makeMessageUpdate(&pages))
}

func (m *Manager) CreateMessage(client bot.Client, channelID snowflake.ID, pages Pages, ephemeral bool) (*discord.Message, error) {
	m.add(&pages)

	return client.Rest().CreateMessage(channelID, m.makeMessageCreate(&pages, ephemeral))
}

func (m *Manager) add(paginator *Pages) {
	m.pagesMu.Lock()
	defer m.pagesMu.Unlock()

	paginator.lastUsed = time.Now()
	m.pages[paginator.ID] = paginator
}

func (m *Manager) remove(paginatorID string) {
	m.pagesMu.Lock()
	defer m.pagesMu.Unlock()
	delete(m.pages, paginatorID)
}

func (m *Manager) OnEvent(event bot.Event) {
	e, ok := event.(*events.ComponentInteractionCreate)
	if !ok {
		return
	}
	customID := e.Data.CustomID()
	if !strings.HasPrefix(customID, m.config.CustomIDPrefix) {
		return
	}
	ids := strings.Split(customID, ":")
	paginatorID, action := ids[1], ids[2]
	paginator, ok := m.pages[paginatorID]
	if !ok {
		if err := e.UpdateMessage(discord.NewMessageUpdateBuilder().ClearContainerComponents().Build()); err != nil {
			e.Client().Logger().Error("Failed to remove components from timed out paginator: ", err)
		}
		return
	}

	if paginator.Creator != 0 && paginator.Creator != e.User().ID {
		if err := e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(m.config.NoPermissionMessage).SetEphemeral(true).Build()); err != nil {
			e.Client().Logger().Error("Failed to send error message: ", err)
		}
		return
	}

	switch action {
	case "first":
		paginator.currentPage = 0

	case "back":
		paginator.currentPage--

	case "stop":
		err := e.UpdateMessage(discord.MessageUpdate{Components: &[]discord.ContainerComponent{}})
		m.remove(paginatorID)
		if err != nil {
			e.Client().Logger().Error("Error updating paginator message: ", err)
		}
		return

	case "next":
		paginator.currentPage++

	case "last":
		paginator.currentPage = paginator.Pages - 1
	}

	if paginator.ExpireMode == ExpireModeAfterLastUsage {
		paginator.lastUsed = time.Now()
	}

	if err := e.UpdateMessage(m.makeMessageUpdate(paginator)); err != nil {
		e.Client().Logger().Error("Error updating paginator message: ", err)
	}
}

func (m *Manager) makeEmbed(paginator *Pages) discord.Embed {
	embedBuilder := discord.NewEmbedBuilder().
		SetFooterText(fmt.Sprintf("Page: %d/%d", paginator.currentPage+1, paginator.Pages)).
		SetColor(m.config.EmbedColor)

	paginator.PageFunc(paginator.currentPage, embedBuilder)
	return embedBuilder.Build()
}

func (m *Manager) makeMessageCreate(pages *Pages, ephemeral bool) discord.MessageCreate {
	var flags discord.MessageFlags
	if ephemeral {
		flags = discord.MessageFlagEphemeral
	}
	return discord.MessageCreate{
		Embeds:     []discord.Embed{m.makeEmbed(pages)},
		Components: []discord.ContainerComponent{m.createComponents(pages)},
		Flags:      flags,
	}
}

func (m *Manager) makeMessageUpdate(pages *Pages) discord.MessageUpdate {
	return discord.MessageUpdate{
		Embeds:     &[]discord.Embed{m.makeEmbed(pages)},
		Components: &[]discord.ContainerComponent{m.createComponents(pages)},
	}
}

func (m *Manager) formatCustomID(paginator *Pages, action string) string {
	return m.config.CustomIDPrefix + ":" + paginator.ID + ":" + action
}

func (m *Manager) createComponents(pages *Pages) discord.ContainerComponent {
	cfg := m.config.ButtonsConfig
	var actionRow discord.ActionRowComponent

	if cfg.First != nil {
		actionRow = actionRow.AddComponents(discord.ButtonComponent{
			Style:    cfg.First.Style,
			Label:    cfg.First.Label,
			Emoji:    &cfg.First.Emoji,
			CustomID: m.formatCustomID(pages, "first"),
			Disabled: pages.currentPage == 0,
		})
	}
	if cfg.Back != nil {
		actionRow = actionRow.AddComponents(discord.ButtonComponent{
			Style:    cfg.Back.Style,
			Label:    cfg.Back.Label,
			Emoji:    &cfg.Back.Emoji,
			CustomID: m.formatCustomID(pages, "back"),
			Disabled: pages.currentPage == 0,
		})
	}

	if cfg.Stop != nil {
		actionRow = actionRow.AddComponents(discord.ButtonComponent{
			Style:    cfg.Stop.Style,
			Label:    cfg.Stop.Label,
			Emoji:    &cfg.Stop.Emoji,
			CustomID: m.formatCustomID(pages, "stop"),
		})
	}

	if cfg.Next != nil {
		actionRow = actionRow.AddComponents(discord.ButtonComponent{
			Style:    cfg.Next.Style,
			Label:    cfg.Next.Label,
			Emoji:    &cfg.Next.Emoji,
			CustomID: m.formatCustomID(pages, "next"),
			Disabled: pages.currentPage == pages.Pages-1,
		})
	}
	if cfg.Last != nil {
		actionRow = actionRow.AddComponents(discord.ButtonComponent{
			Style:    cfg.Last.Style,
			Label:    cfg.Last.Label,
			Emoji:    &cfg.Last.Emoji,
			CustomID: m.formatCustomID(pages, "last"),
			Disabled: pages.currentPage == pages.Pages-1,
		})
	}

	return actionRow
}
