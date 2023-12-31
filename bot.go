package main

import (
	// Hashing and Associated Libraries
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	// Database Libraries
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	// Standard Libraries
	"log"
	"os"
	"sync"
	"time"

	// Third-Party Libraries
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

func main() {

	// Use third-party library to load environment variables
	err := godotenv.Load()
	if err != nil {
		//	log.Fatal("Error Loading Environment Variables")
	}

	// Bot Info
	Token := os.Getenv("TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	// Attempt to Connect to Discord API for Bot using Credentials
	discord_agent, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error Creating Discord Bot: %v", err)
	}

	// Build Cron Job Scheduler
	c := cron.New()

	// Check Every 15 Minutes
	c.AddFunc("0/15 0/23 * * *", func() {
		var con sync.WaitGroup
		// Allow for Multiple URLS to be Read from the Feed
		feedURLs := []string{
			"http://news.fullerton.edu/engineering-and-computer-science/feed",
			"http://news.fullerton.edu/business-and-economics/feed",
			"http://news.fullerton.edu/natural-sciences-and-mathematics/feed",
			"http://news.fullerton.edu/all-news/feed",
		}

		// Loop Over ignoring index of posts
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

func parseRSS(session *discordgo.Session, channelID string, feedURL string) {

	// RSS Feed Parser Object for Go
	pRss := gofeed.NewParser()

	rssFeed, err := pRss.ParseURL(feedURL)
	// Check for Errors when Parsing RSS Feeds
	if err != nil {
		log.Printf("Error Fetching: %v", err)
	}

	if rssFeed == nil {
		log.Println("No RSS Feed Found")
	}

	// Loop iterates over the last 15 posts; 0 is the most recent and 15 is the oldest.
	// This is to prevent spamming the channel in Reverse Chronological Order
	for i := 15; i != 0; i-- {

		// Open Database Connection
		db, err := sql.Open("sqlite3", "./var/announcements.db")

		if err != nil {
			log.Fatal("Error Opening Database Connection to announcements.db" + err.Error())
		}

		var item = rssFeed.Items[i]
		// Needed to Pull Images Specifically from the RSS Feeds from
		mediaContent := item.Extensions["media"]["content"]
		mediaURL := mediaContent[0].Attrs["url"]

		combinedHash := generateHash(item.Title, item.Link, *item.PublishedParsed)

		// Count Variable for query
		var count int
		// Check if the Hash is in the Database
		db_err := db.QueryRow("SELECT COUNT(*) FROM announcements WHERE combined_hash = ?", combinedHash).Scan(&count)

		if db_err != nil {
			log.Printf("Error Querying Database: %v", err)
		} else if count == 0 {
			// Insert the Hash into the Database and Create the Posts

			// Inserting the Hash
			db.Exec("INSERT INTO announcements (combined_hash) VALUES (?)", combinedHash)

			// Create the Embedded Message
			embed := &discordgo.MessageEmbed{
				Title:       item.Title,
				URL:         item.Link,
				Description: item.Description,
				Image:       &discordgo.MessageEmbedImage{URL: mediaURL},
			}

			// Check if the Image is Nil
			if item.Image != nil {
				embed.Image = &discordgo.MessageEmbedImage{URL: item.Image.URL}
			} else if rssFeed.Image != nil {
				embed.Image = &discordgo.MessageEmbedImage{URL: rssFeed.Image.URL}
			}

			// Send the Embedded Message
			_, err = session.ChannelMessageSendEmbed(channelID, embed)
			if err != nil {
				log.Printf("Error sending embedded message: %v", err)
			}
		} else {
			log.Println("No New Posts Found")
		}

		// Close Database Connection
		db.Close()
	}
}

func generateHash(title string, link string, date time.Time) string {

	// Combine the Title, Link, and URL i
	combinedString := fmt.Sprintf("%s%s%s", title, link, date.String())

	// Hash the Combined String
	hash := sha256.New()
	hash.Write([]byte(combinedString))
	combinedHash := hex.EncodeToString(hash.Sum(nil))
	return combinedHash
}
