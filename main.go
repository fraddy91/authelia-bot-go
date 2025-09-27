package main

import (
	"log"
	"os"

	"github.com/fraddy91/authelia-bot-go/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	bot.EnsureStartupFiles()
	// Load token from environment
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}
	log.Println("Bot starting...")
	log.Println("BOT_TOKEN =", os.Getenv("BOT_TOKEN"))

	// Init bot
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	// Add commands info
	commands := []tgbotapi.BotCommand{
		{Command: "chatid", Description: "Show chatId"},
		{Command: "whoami", Description: "Show user registration status"},
		{Command: "mode", Description: "Change notification mode"},
		{Command: "register", Description: "Request registration"},
		{Command: "unregister", Description: "Remove registered user"},
		{Command: "approve", Description: "Approve pending user"},
		{Command: "deny", Description: "Deny pending user"},
		{Command: "unignore", Description: "Remove user from ignore list"},
		{Command: "notify", Description: "Toggle admin alerts"},
		{Command: "pendings", Description: "List pending registrations"},
		{Command: "health", Description: "Show current status"},
		{Command: "users", Description: "Show users"},
		{Command: "menu", Description: "Show admin panel"},
	}

	setCmd := tgbotapi.NewSetMyCommands(commands...)
	if _, err := botAPI.Request(setCmd); err != nil {
		log.Println("Failed to set command menu:", err)
	}

	// Start file watcher in background
	go bot.WatchNotifications(botAPI)

	// Start updates channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	// Main loop
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "chatid":
				bot.HandleChatID(botAPI, update)
			case "whoami":
				bot.HandleWhoAmI(botAPI, update)
			case "mode":
				bot.HandleMode(botAPI, update)
			case "pendings":
				bot.HandlePendings(botAPI, update)
			case "health":
				bot.HandleHealth(botAPI, update)
			case "register":
				bot.HandleRegister(botAPI, update)
			case "unregister":
				bot.HandleUnregister(botAPI, update)
			case "unignore":
				bot.HandleUnignore(botAPI, update)
			case "notify":
				bot.HandleNotify(botAPI, update)
			case "users":
				bot.HandleUsers(botAPI, update)
			case "menu":
				bot.HandleMenu(botAPI, update)
			}
		}

		if update.CallbackQuery != nil {
			bot.HandleCallback(botAPI, update)
		}
	}
}
