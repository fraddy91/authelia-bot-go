package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------------- Command Handlers ----------------------

func HandleChatID(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🆔 Your chat ID is: %d", chatID))
	bot.Send(msg)
}

func HandleWhoAmI(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	users := loadUsers()
	for email, u := range users {
		if u.ChatID == chatID {
			msg := fmt.Sprintf("📧 Email: %s\n🔔 Mode: %s\n👑 Admin: %v\n🔔 Notify: %v",
				email, u.Mode, u.Admin, u.Notify)
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not registered."))
}

func HandleMode(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 1 || (args[0] != "full" && args[0] != "short" && args[0] != "code") {
		bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Usage: /mode <full|short|code>"))
		return
	}
	newMode := args[0]
	users := loadUsers()
	for email, u := range users {
		if u.ChatID == chatID {
			u.Mode = newMode
			users[email] = u
			saveUsers(users)
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Updated mode to %s", newMode)))
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not registered."))
}

func HandlePendings(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if !isAdmin(chatID) {
		bot.Send(tgbotapi.NewMessage(chatID, "🚫 Admins only."))
		return
	}
	pending := loadPending()
	if len(pending) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "✅ No pending registrations."))
		return
	}
	count := 0
	for pid, req := range pending {
		if count >= 10 {
			break
		}
		email := req.Email
		mode := req.Mode
		cid := req.ChatID
		msg := tgbotapi.NewMessage(chatID,
			fmt.Sprintf("⏳ Pending:\n📧 %s\n🆔 chat %s\n🔔 mode %s", email, pid, mode))
		msg.ReplyMarkup = buildApprovalKeyboard(cid)
		bot.Send(msg)
		count++
	}
}

func HandleHealth(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if !isAdmin(chatID) {
		bot.Send(tgbotapi.NewMessage(chatID, "🚫 Admins only."))
		return
	}
	users := loadUsers()
	pending := loadPending()
	ignore := loadIgnore()
	uptime := int(time.Since(StartTime).Seconds())
	msg := fmt.Sprintf(
		"🩺 Healthcheck\n⏱️ Uptime: %ds\n👥 Users: %d\n⏳ Pending: %d\n🚫 Ignored: %d\n📨 Last notification: %v",
		uptime, len(users), len(pending), len(ignore), LastProcessed.UTC().Format("2006-01-02 15:04:05 UTC"),
	)

	bot.Send(tgbotapi.NewMessage(chatID, msg))
}

func SendAdminMenu(bot *tgbotapi.BotAPI, chatID int64) {
	menu := tgbotapi.NewMessage(chatID, "🛠️ Admin Menu:")
	menu.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/pendings"),
			tgbotapi.NewKeyboardButton("/users"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/notify on"),
			tgbotapi.NewKeyboardButton("/notify off"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/unignore"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/health"),
		),
	)
	bot.Send(menu)
}

func HandleMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	users := loadUsers()

	for _, user := range users {
		if user.ChatID == chatID && user.Admin {
			SendAdminMenu(bot, chatID)
			return
		}
	}

	bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not authorized to view the admin menu."))
}

// ---------------------- Notify Handler ----------------------

func HandleNotify(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 1 || (args[0] != "on" && args[0] != "off") {
		bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Usage: /notify on|off"))
		return
	}

	users := loadUsers()
	for email, user := range users {
		if user.ChatID == chatID {
			user.Notify = (args[0] == "on")
			users[email] = user
			saveUsers(users)
			status := "enabled"
			if !user.Notify {
				status = "disabled"
			}
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("🔔 Notifications %s for your account.", status)))
			return
		}
	}

	bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not a registered user."))
}

// ---------------------- Unignore Handler ----------------------

func HandleUnignore(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 1 {
		bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Usage: /unignore <chat_id>"))
		return
	}

	targetID := args[0]
	ignore := loadIgnore()
	if _, exists := ignore[targetID]; exists {
		delete(ignore, targetID)
		saveIgnore(ignore)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Chat %s removed from ignore list.", targetID)))
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("ℹ️ Chat %s is not in ignore list.", targetID)))
	}
}

// ---------------------- Callback Handler ----------------------

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cb := update.CallbackQuery
	adminChat := cb.Message.Chat.ID

	if !isAdmin(adminChat) {
		bot.Request(tgbotapi.NewCallback(cb.ID, "🚫 Admins only."))
		return
	}

	parts := strings.Split(cb.Data, ":")
	if len(parts) != 2 {
		bot.Request(tgbotapi.NewCallback(cb.ID, "⚠️ Invalid action."))
		return
	}
	action, chatIDStr := parts[0], parts[1]

	pending := loadPending()
	req, ok := pending[chatIDStr]
	if !ok {
		bot.Request(tgbotapi.NewCallback(cb.ID, "❌ No such pending request."))
		return
	}

	switch action {
	case "approve":
		users := loadUsers()
		users[req.Email] = User{
			ChatID: req.ChatID,
			Mode:   req.Mode,
			Admin:  false,
			Notify: false,
		}
		saveUsers(users)
		delete(pending, chatIDStr)
		savePending(pending)

		bot.Request(tgbotapi.NewCallback(cb.ID, "✅ Approved"))
		edit := tgbotapi.NewEditMessageText(adminChat, cb.Message.MessageID,
			fmt.Sprintf("✅ Approved %s (chat %s)", req.Email, chatIDStr))
		bot.Send(edit)

		notify := tgbotapi.NewMessage(req.ChatID,
			fmt.Sprintf("✅ Your registration for `%s` has been approved.\nMode: `%s`", req.Email, req.Mode))
		notify.ParseMode = "Markdown"
		bot.Send(notify)

	case "deny":
		delete(pending, chatIDStr)
		savePending(pending)
		ignore := loadIgnore()
		ignore[chatIDStr] = true
		saveIgnore(ignore)

		bot.Request(tgbotapi.NewCallback(cb.ID, "🚫 Denied"))
		edit := tgbotapi.NewEditMessageText(adminChat, cb.Message.MessageID,
			fmt.Sprintf("🚫 Denied and blocked chat %s", chatIDStr))
		bot.Send(edit)

		notify := tgbotapi.NewMessage(req.ChatID,
			fmt.Sprintf("🚫 Your registration for `%s` has been denied.", req.Email))
		notify.ParseMode = "Markdown"
		bot.Send(notify)
	}
}

func HandleUsers(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	users := loadUsers()

	// Check if sender is an admin
	isAdmin := false
	for _, user := range users {
		if user.ChatID == chatID && user.Admin {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not authorized to view the user list."))
		return
	}

	if len(users) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "📭 No registered users found."))
		return
	}

	var lines []string
	for email, user := range users {
		adminFlag := ""
		if user.Admin {
			adminFlag = "🛡️"
		}
		notifyFlag := "🔔"
		if !user.Admin || !user.Notify {
			notifyFlag = "🔕"
		}
		lines = append(lines, fmt.Sprintf("%s %s\nChatID: %d\nMode: %s %s\n", adminFlag, email, user.ChatID, user.Mode, notifyFlag))
	}

	msg := tgbotapi.NewMessage(chatID, "📋 Registered Users:\n\n"+strings.Join(lines, "\n"))
	bot.Send(msg)
}

func HandleRegister(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	args := strings.Fields(update.Message.CommandArguments())
	if len(args) != 2 {
		bot.Send(tgbotapi.NewMessage(chatID, "⚠️ Usage: /register <email> <mode>"))
		return
	}
	email := args[0]
	mode := args[1]
	if mode != "full" && mode != "short" && mode != "code" {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Invalid mode. Choose: full, short, code"))
		return
	}

	// Check ignore list
	ignore := loadIgnore()
	if ignore[fmt.Sprint(chatID)] {
		bot.Send(tgbotapi.NewMessage(chatID, "🚫 You are blocked from registering."))
		return
	}

	// Add to pending
	pending := loadPending()
	pending[fmt.Sprint(chatID)] = PendingRequest{
		Email:  email,
		Mode:   mode,
		ChatID: chatID,
	}
	savePending(pending)

	bot.Send(tgbotapi.NewMessage(chatID, "⏳ Registration request submitted. Waiting for admin approval."))

	// Notify admins
	admins := loadUsers()
	for _, admin := range admins {
		if admin.Admin && admin.Notify {
			msg := tgbotapi.NewMessage(admin.ChatID,
				fmt.Sprintf("📥 New registration request\n📧 %s\n🆔 chat %d\n🔔 mode %s",
					email, chatID, mode))
			msg.ReplyMarkup = buildApprovalKeyboard(chatID)
			bot.Send(msg)
		}
	}
}

func HandleUnregister(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	if !AllowRegistration {
		bot.Send(tgbotapi.NewMessage(chatID, "🚫 Registration is disabled on this bot."))
		return
	}
	users := loadUsers()
	for email, u := range users {
		if u.ChatID == chatID {
			if u.Admin {
				bot.Send(tgbotapi.NewMessage(chatID, "🚫 Admin accounts cannot unregister via command."))
				return
			}
			delete(users, email)
			saveUsers(users)
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("🗑️ Unregistered %s", email)))
			return
		}
	}
	bot.Send(tgbotapi.NewMessage(chatID, "❌ You are not registered."))
}

// ---------------------- Inline Buttons ----------------------

func buildApprovalKeyboard(chatID int64) tgbotapi.InlineKeyboardMarkup {
	approve := tgbotapi.NewInlineKeyboardButtonData("✅ Approve", fmt.Sprintf("approve:%d", chatID))
	deny := tgbotapi.NewInlineKeyboardButtonData("❌ Deny", fmt.Sprintf("deny:%d", chatID))
	row := tgbotapi.NewInlineKeyboardRow(approve, deny)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}
