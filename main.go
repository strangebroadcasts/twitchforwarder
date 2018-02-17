package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/pkg/browser"

	"github.com/strangebroadcasts/twitchforwarder/discord"
)

// HookURL holds the URL for the Discord webhook
var HookURL string

func main() {
	// twitchChannel := flag.String("channel", "", "the Twitch channel to forward messages from")
	nickname := flag.String("nick", "", "Your Twitch username")
	channel := flag.String("channel", "", "The Twitch channel to forward messages from")
	oauthToken := flag.String("oauth", "", "OAuth token for Twitch")
	flag.StringVar(&HookURL, "hook", "", "the URL (with ID and token) for the Discord webhook")
	flag.Parse()

	if *nickname == "" {
		log.Fatal("Twitch username must be set with -nick")
	}

	if *channel == "" {
		log.Fatal("Twitch channel must be set with -channel")
	}

	// Twitch requires nicknames in lowercase when connecting via IRC:
	*nickname = strings.ToLower(*nickname)

	if HookURL == "" {
		log.Fatal("Discord hook URL must be set with -hook")
	}

	// If OAuth token is missing, prompt the user to generate one:
	if *oauthToken == "" {
		fmt.Println("No OAuth token. Opening https://twitchapps.com/tmi/...")
		err := browser.OpenURL("https://twitchapps.com/tmi/")
		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Println("Authorize the app, and copy the token here, including \"oauth:\".")
		fmt.Print("OAuth Token: ")

		var newToken string
		fmt.Scan(&newToken)
		*oauthToken = newToken
		fmt.Println(*oauthToken)
	}

	client := twitch.NewClient(*nickname, *oauthToken)
	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		postMessage(user.Username, message.Text)
	})

	client.Join(*channel)
	go func() {
		err := client.Connect()
		if err != nil {
			log.Fatal("Error:", err)
		}
	}()

	log.Print("Now running - Ctrl+C to shut down.")

	// Wait until CTRL+C:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Shutting down.")

	client.Disconnect()
	return
}

func postMessage(username, content string) error {
	exec := discord.WebhookExecution{Content: content, Username: "[TWITCH] " + username}
	json, err := json.Marshal(exec)
	if err != nil {
		return err
	}

	jsonReader := bytes.NewReader(json)

	resp, err := http.Post(HookURL, "application/json", jsonReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		explanation := fmt.Sprint("Got error ", resp.StatusCode, resp.Status, "from Discord")
		return errors.New(explanation)
	}

	return nil
}
