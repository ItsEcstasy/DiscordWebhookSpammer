package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	WebhookLink string `json:"webhookLink"`
}

func main() {
	config, err := loadConfig("settings.json")
	if err != nil {
		fmt.Println("[!] Error loading configuration:", err)
		return
	}

	message, err := ReadMSG("message.txt")
	if err != nil {
		fmt.Println("[!] Error reading message from file:", err)
		return
	}

	session, err := discordgo.New(config.WebhookLink)
	if err != nil {
		fmt.Println("[!] Error creating Discord session:", err)
		return
	}

	webhookID, webhookToken, err := parseWebhook(config.WebhookLink)
	if err != nil {
		fmt.Println("[!] Error parsing webhook link:", err)
		return
	}

	session.WebhookExecute(webhookID, webhookToken, false, &discordgo.WebhookParams{
		Content: message,
	})

	fmt.Println("[+] Starting Spam.")

	for {
		session.WebhookExecute(webhookID, webhookToken, false, &discordgo.WebhookParams{
			Content: message,
		})

		time.Sleep(100 * time.Millisecond)
	}
}

func loadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ReadMSG(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

func parseWebhook(webhookLink string) (webhookID string, webhookToken string, err error) {
	parts := strings.Split(webhookLink, "/")
	if len(parts) < 7 {
		return "", "", fmt.Errorf("[!] Invalid webhook link")
	}

	webhookID = parts[5]
	webhookToken = parts[6]

	return webhookID, webhookToken, nil
}
