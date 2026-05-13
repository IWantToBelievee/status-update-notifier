package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nShutting down...")
		cancel()
	}()

	parser := gofeed.NewParser()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := loadConfig()
	if cfg.url == "" || cfg.flag == "" {
		fmt.Println("URL and FLAG environment variables must be set.")
		cancel()
	}

	go checkStatus(ctx, parser, cfg)

	<-ctx.Done()
}

type config struct {
	url            string
	flag           string
	check_interval time.Duration
	status_index   int
}

func loadConfig() config {
	checkInterval, err := strconv.Atoi(os.Getenv("CHECK_INTERVAL"))
	if err != nil {
		checkInterval = 300 // Default to 5 minutes
	}

	statusIndex, err := strconv.Atoi(os.Getenv("STATUS_INDEX"))
	if err != nil {
		statusIndex = 0 // Default to the first element
	}

	return config{
		url:            os.Getenv("URL"),
		flag:           os.Getenv("FLAG"),
		check_interval: time.Duration(checkInterval) * time.Second,
		status_index:   statusIndex,
	}
}

func checkStatus(ctx context.Context, parser *gofeed.Parser, cfg config) {
	// Use a map to track seen items and avoid duplicate notifications
	seenItems := make(map[string]bool)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Fetch and parse the RSS feed
			feed, err := parser.ParseURL(cfg.url)
			if err != nil {
				fmt.Printf("[%s] Error fetching feed: %v\n", time.Now().Format(time.TimeOnly), err)

				select {
				case <-time.After(5 * time.Second):
				case <-ctx.Done():
					return
				}
				continue
			}

			if len(feed.Items) > 0 {
				// Check for items containing the specified flag and notify the user if a new status is found
				for _, item := range feed.Items {
					if strings.Contains(item.Title, cfg.flag) && !seenItems[item.Title] {
						seenItems[item.Title] = true

						// Extract the status text based on the specified index
						splitted := strings.Split(item.Title, " ")
						statusText := splitted[cfg.status_index]

						notifyUser(fmt.Sprintf("Status changed: %s", statusText))
					}
				}
			}

			// Wait for the specified check interval before checking again
			select {
			case <-time.After(cfg.check_interval):
			case <-ctx.Done():
				return
			}
		}
	}
}

// notifyUser sends a desktop notification with the given message
func notifyUser(message string) {
	cmd := exec.Command("notify-send", "Status Update", message)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[%s] Error sending notification: %v\n", time.Now().Format(time.TimeOnly), err)
	} else {
		fmt.Printf("[%s] Notification sent: %s\n", time.Now().Format(time.TimeOnly), message)
	}
}
