[![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/paginator.svg)](https://pkg.go.dev/github.com/disgoorg/disgo)
[![Go Report](https://goreportcard.com/badge/github.com/disgoorg/paginator)](https://goreportcard.com/report/github.com/disgoorg/paginator)
[![Go Version](https://img.shields.io/github/go-mod/go-version/disgoorg/paginator)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/paginator/blob/master/LICENSE)
[![Paginator Version](https://img.shields.io/github/v/tag/disgoorg/paginator?label=release)](https://github.com/disgoorg/paginator/releases/latest)
[![DisGo Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/TewhTfDpvW)

# Paginator

Paginator is a simple embed pagination library for [DisGo](https://github.com/disgoorg/disgo) using buttons. It supports both interactions and normal messages.

## Getting Started

### Installation

```bash
$ go get github.com/disgoorg/paginator
```

### Usage

Create a new paginator and add it as `EventHandler` to your Client.
```go
manager := paginator.New()
client, err := disgo.New(token,
    bot.WithDefaultGateway(),
    bot.WithEventListeners(manager),
)
```

#### Interactions

```go
// your data to paginate through
pData := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}

// create a new paginator.Pages
err := manager.Create(event.Respond, paginator.Pages{
    // A unique ID for this paginator
    ID: event.ID().String(), 
	// This is the function that will be called to create the embed for each page when the page is displayed
    PageFunc: func(page int, embed *discord.EmbedBuilder) {
        embed.SetTitle("Data")
        description := ""
        for i := 0; i < 5; i++ {
            if page*5+i >= len(pData) {
                break
            }
            description += pData[page*5+i] + "\n"
        }
        embed.SetDescription(description)
    },
	// The total number of pages
    Pages:      int(math.Ceil(float64(len(pData)) / 5)),
	// Optional: If the paginator should only be accessible by the user who created it
    Creator:    event.User().ID,
	// Optional: If the paginator should be deleted after x time after the last interaction
    ExpireMode: paginator.ExpireModeAfterLastUsage,
}, false)
```

#### Normal Messages

```go
// your data to paginate through
pData := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
channelID := 123454545
// create a new paginator.Pages
err := manager.CreateMessage(event.Client(), channelID, paginator.Pages{
    // A unique ID for this paginator
    ID: event.Message.ID.String(), 
	// This is the function that will be called to create the embed for each page when the page is displayed
    PageFunc: func(page int, embed *discord.EmbedBuilder) {
        embed.SetTitle("Data")
        description := ""
        for i := 0; i < 5; i++ {
            if page*5+i >= len(pData) {
                break
            }
            description += pData[page*5+i] + "\n"
        }
        embed.SetDescription(description)
    },
	// The total number of pages
    Pages:      int(math.Ceil(float64(len(pData)) / 5)),
	// Optional: If the paginator should only be accessible by the user who created it
    Creator:    event.User().ID,
	// Optional: If the paginator should be deleted after x time after the last interaction
    ExpireMode: paginator.ExpireModeAfterLastUsage,
}, false)
```

## Documentation

Documentation is wip and can be found under

* [![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/paginator.svg)](https://pkg.go.dev/github.com/disgoorg/paginator)

## Examples

You can find examples [here](https://github.com/disgoorg/paginator/tree/master/_example)

## Troubleshooting

For help feel free to open an issue or reach out on [Discord](https://discord.gg/TewhTfDpvW)

## Contributing

Contributions are welcomed but for bigger changes we recommend first reaching out via [Discord](https://discord.gg/TewhTfDpvW) or create an issue to discuss your problems, intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/paginator/blob/master/LICENSE). See LICENSE for more information.
