package main

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

func parseRSS(session *discordgo.Session, channelID string, feedURL string) {
	pRss := gofeed.NewParser()
	rssFeed, err := pRss.ParseURL(feedURL)
	// Check for Errors when Parsing RSS Feeds
	if err != nil {
		log.Printf("Error Fetching: %v", err)
	}

	if rssFeed == nil {
		log.Println("No RSS Feed Found")
	}

	for i := 0; i < 5; i++ {
		var item = rssFeed.Items[i]
		// Needed to Pull Images Specifically from the RSS Feeds from
		mediaContent := item.Extensions["media"]["content"]
		mediaURL := mediaContent[0].Attrs["url"]

		embed := &discordgo.MessageEmbed{
			Title:       item.Title,
			URL:         item.Link,
			Description: item.Description,
			Image:       &discordgo.MessageEmbedImage{URL: mediaURL},
		}

		if item.Image != nil {
			embed.Image = &discordgo.MessageEmbedImage{URL: item.Image.URL}
		} else if rssFeed.Image != nil {
			embed.Image = &discordgo.MessageEmbedImage{URL: rssFeed.Image.URL}
		}

		_, err = session.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Printf("Error sending embed: %v", err)
		}
	}
}

func main() {

	// Bot Info
	Token := "MTE0MzY3NjU4MjIzODYzODIyMA.G4J3ec.JzoBhY3ehwbkBCAl3CwJS7lf0HnGke1EsqwUe8"
	channelID := "1145180678468665475"

	discord_agent, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error Creating Discord Bot: %v", err)
	}

	// Build Cron Job Scheduler
	c := cron.New()

	// Check Every 15 Minutes
	c.AddFunc("0/15 0/23 * * *", func() { // Reset to " * * * * * " for Testing
		var con sync.WaitGroup
		// Allow for Multiple URLS to be Read from the Feed
		feedURLs := []string{
			"http://news.fullerton.edu/engineering-and-computer-science/feed",
			"http://news.fullerton.edu/business-and-economics/feed",
		}

		// Loop Over ignoring index, hence _,
		for _, urls := range feedURLs {
			con.Add(1)
			go func(url string) {
				defer con.Done()
				parseRSS(discord_agent, channelID, url)
			}(urls)
		}

		con.Wait()
	})

	c.Start()
	select {}
}
