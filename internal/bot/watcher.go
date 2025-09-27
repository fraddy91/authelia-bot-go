package bot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------------- Parsing helpers ----------------------

func extractEmail(text string) string {
	re := regexp.MustCompile(`Recipient:\s*\{.*?([\w\.-]+@[\w\.-]+).*?\}`) //or generic? MustCompile(`(?i)([a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,})`)
	m := re.FindStringSubmatch(text)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func extractCode(text string) string {
	re := regexp.MustCompile(`-{80}\s*\n([A-Z0-9\.]+)\s*\n-{80}`)
	m := re.FindStringSubmatch(text)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func extractLinks(text string) []string {
	re := regexp.MustCompile(`https://[^\s]+/revoke/one-time-code\?id=\w+`)
	all := re.FindAllString(text, -1)
	seen := map[string]bool{}
	out := []string{}
	for _, l := range all {
		if !seen[l] {
			seen[l] = true
			out = append(out, l)
		}
	}
	return out
}

func formatMessage(text, mode string) string {
	switch mode {
	case "full":
		if len(text) > 4096 {
			return text[:4096]
		}
		return text
	case "short":
		code := extractCode(text)
		links := extractLinks(text)
		return fmt.Sprintf("🔐 Code: %s\n🔗 Links:\n%s", code, strings.Join(links, "\n"))
	case "code":
		code := extractCode(text)
		return fmt.Sprintf("🔐 Code: %s", code)
	default:
		return "⚠️ Unknown mode"
	}
}

// ---------------------- Watcher ----------------------

func WatchNotifications(bot *tgbotapi.BotAPI) {
	log.Println("Watcher starting...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Failed to create watcher:", err)
	}
	defer watcher.Close()

	// Ensure notifications dir exists
	os.MkdirAll("notifications", 0755)

	path := "notifications/notification.txt"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("notification.txt not found, creating empty file...")
		os.WriteFile(path, []byte(""), 0644)
	}

	err = watcher.Add("notifications")
	if err != nil {
		log.Fatal("Failed to watch notifications dir:", err)
	}
	log.Println("👀 Watching notifications.txt for changes...")

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write &&
				strings.HasSuffix(event.Name, "notification.txt") {
				log.Printf("📨 File notifications.txt changed")

				LastProcessed = time.Now()

				content, err := os.ReadFile(path)
				if err != nil {
					log.Println("Error reading notification.txt:", err)
					continue
				}

				log.Printf("📨 Read %d bytes from notifications.txt", len(content))
				log.Printf("📄 Content: %s", string(content))

				email := extractEmail(string(content))
				if email == "" {
					log.Println("No email found in notification")
					continue
				}

				users := loadUsers()
				user, ok := users[email]
				if !ok {
					log.Printf("No user found for email: %s", email)
					continue
				}

				msg := tgbotapi.NewMessage(user.ChatID,
					formatMessage(string(content), user.Mode))
				bot.Send(msg)

				log.Printf("Notification sent to %s (chat %d)", email, user.ChatID)
			}

		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
