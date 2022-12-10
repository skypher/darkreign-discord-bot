package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

// DiscordBot represents a Discord bot.
type DiscordBot struct {
	URL       string
	LastTitle string
	LastBody  string
	Session   *discordgo.Session
}

// Start starts the Discord bot.
func (bot *DiscordBot) Start() {
	discord, err := discordgo.New()
	if err != nil {
		log.Fatal(err)
	}

	discord.Token = "YOUR_BOT_TOKEN"
	discord.LogLevel = discordgo.LogInformational

	bot.Session = discord

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				bot.ScrapeWebPage()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// ScrapeWebPage scrapes the specified URL for the page title and first paragraph.
func (bot *DiscordBot) ScrapeWebPage() {
	doc, err := goquery.NewDocument(bot.URL)
	if err != nil {
		log.Fatal(err)
	}

	title := bot.GetTitle(doc)
	paragraph := bot.GetParagraph(doc)

	changes := bot.CheckForChanges(title, paragraph)

	if changes {
		bot.SendMessage(title, paragraph)
		bot.UpdateLastScrape(title, paragraph)
	}
}

// GetTitle gets the page title from the specified goquery.Document.
func (bot *DiscordBot) GetTitle(doc *goquery.Document) string {
	title := doc.Find("head > title").First().Text()
	return title
}

// GetParagraph gets the first paragraph from the specified goquery.Document.
func (bot *DiscordBot) GetParagraph(doc *goquery.Document) string {
	paragraph := doc.Find("p").First().Text()
	return paragraph
}

// CheckForChanges checks if the page title or first paragraph have changed since the last scrape.
func (bot *DiscordBot) CheckForChanges(title, body string) bool {
	if title != bot.LastTitle || body != bot.LastBody {
		return true
	}

	return false
}

// SendMessage sends a message to the Discord channel with the specified page title and first paragraph.
func (bot *DiscordBot) SendMessage(title, body string) {
	channelID, err := bot.Session.State.ChannelID("YOUR_CHANNEL_ID")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.Session.ChannelMessageSend(channelID, fmt.Sprintf("**%s**\n%s", title, body))
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateLastScrape updates the last scraped title and first paragraph.
func (bot *DiscordBot) UpdateLastScrape(title, body string) {
	bot.LastTitle = title
	bot.LastBody = body
}

func main() {
	bot := &DiscordBot{
		URL: "https://www.example.com",
	}

	bot.Start()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Session.Close()
}
