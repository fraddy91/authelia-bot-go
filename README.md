![GitHub Release](https://img.shields.io/github/v/release/fraddy91/authelia-bot-go)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/fraddy/authelia-bot?logo=docker&sort=date)
![Docker Image Size](https://img.shields.io/docker/image-size/fraddy/authelia-bot)

# Authelia Telegram Bot (Go)

A lightweight, modular, admin-moderated Telegram bot for delivering notifications from Authelia's Notification File Provider with manual registration approval, ignore lists, and abuse-resistant workflows. Built in Go, Dockerized and ready for CI/CD.
<br>Memory footprint is just 5 mb in standby.

I made this bot just to get OTP from Authelia, because there are no any solutions to get notifications anywhere else than email.
<br>And also for fun, and to learn Go a little. 😸

---

This bot acts as a secure bridge between Authelia and Telegram, delivering sensitive notifications (such as login attempts, 2FA codes, or access alerts) directly to approved users.
<br>The workflow:
1. 	Authelia writes a notification (e.g. login attempt, code, alert) to `notification.txt` via Notification File Provider.
2. 	The bot watches this file for changes.
3. 	When updated, it extracts the target email and message content.
4. 	If the email matches a registered user, the bot sends the message to their Telegram chat.
This ensures:
• 	✅ Only approved users receive sensitive notifications
• 	✅ Messages are delivered instantly and securely
• 	✅ Admins retain full control over who gets notified


## 🚀 Features

- ✅ Manual registration with admin approval
- ✅ Inline buttons for approve/deny
- ✅ Ignore list support
- ✅ Notification delivery via `notification.txt`
- ✅ Admin opt-in alerts for new registrations
- ✅ Environment-driven toggles (`ALLOW_REGISTRATION`)
- ✅ Docker-ready with volume mounts
- ✅ Healthcheck and pprof support
- ✅ Hardened file watcher with fallback logic

---

## 🗨️ Commands

- /register \<email> \<mode> — Request registration. Mode must be one of: full, short, or code.
- /unignore \<chat_id> — Remove a user from the ignore list.
- /notify on or /notify off — Enable or disable admin notifications for new registration requests.
- /pendings — Show all pending registrations.
- /users — List all registered users.
- /chat Show your current Telegram chat ID.
- /whoami Show your registration info (email, mode, admin status, notify status).
- /mode \<new_mode> Change your delivery mode (full, short, code).
- /menu Show a persistent button menu (admins only).
- /health Show last notification timestamp and bot status.
- /unregister Remove yourself from the bot’s user list.

### Callbacks

- /approve \<chat_id> — Approve a pending user by their Telegram chat ID.
- /deny \<chat_id> — Deny a pending user by chat ID.

## 🔨 Configuration

### Environment Variables (set via .env or Docker/Docker-Compose environments)

- BOT_TOKEN — your Telegram bot token (required).
- ALLOW_REGISTRATION — set to true or false to enable or disable user registration.
- ENABLE_PPROF — optional; set to true to enable pprof debugging on port 6060.

### JSON Files

- config/users.json
Stores registered users with their email, chat ID, mode, admin status, and notification preference.
- config/pending.json
Stores pending registration requests keyed by chat ID.
- config/ignore.json
Stores blocked users keyed by chat ID.
- notifications/notifications.txt
Should be mapped to Authelia's notifications.txt. When updated, the bot reads this file and delivers its contents to the matching user based on email.

## 👤 Usage

### Approval & Denial Workflow
When a user sends /register \<email> \<mode>, the bot places them in a pending queue.
<br>Admins receive a notification with inline buttons:

<br>New registration request:
<br>Email: user@example.com
<br>Mode: full
<br>ChatID: 123456789
[✅ Approve] [❌ Deny]

- Tapping Approve registers the user and stores their info in users.json.
- Tapping Deny removes them from the pending queue.
- These buttons use Telegram callback queries, so responses are instant and don’t clutter the chat.


## 🛠️ Setup

### 🐳 Docker Deployment
You can run the bot in Docker using the following setup.
<br>Folder structure
```
authelia-bot-go/
├── config/                # Stores users.json, pending.json, ignore.json
├── notifications/         # Authelia writes notification.txt here
├── Dockerfile
├── docker-compose.yml
```

```yml
Sample docker-compose.yml
version: '3.8'
services:
  notifier:
    image: fraddy/authelia-bot:latest
    container_name: notifier
    environment:
      - BOT_TOKEN=your-telegram-bot-token
      - ALLOW_REGISTRATION=true
    volumes:
      - ./config:/app/config # Persistent folder
      - {Authelia's folder with notification.txt}/:/app/notifications/ # Map to the Authelia's data folder with notification.txt
    ports:
      - "6060:6060"  # optional: for pprof debugging
    restart: unless-stopped
```

### First steps

You should configure BOT_TOKEN variable, add users to users.json and you're good to go.

users.json example:
```yml
{
  "alice@example.com": {
    "chat_id": 1234567890,
    "mode": "code",
    "admin": false,
    "notify": false
  }
}
```

You can ask chatId from bot, using /chatId command, or use any other way.

### Registration

If you've turned on registrations, then you'll get pending registrations in pending.json file:

```yml
{
  "1234567892": {
    "email": "bob@example.com",
    "mode": "full",
    "chat_id": 1234567892
  }
}
```

Any admin that has notify setting can approve/deny registration using callback or manually get /pending list for this.

### Blacklist

You can manage blacklist using /unignore command

```yml
{
  "1234567891": true
}
```

## Result examples

<details>
  <summary> Generic Commands screennshots </summary>

### Generic Commands example
![Example 1.](/assets/example1.png)
### Perdings, mode, whoami<br>
![Example 4.](/assets/example4.png)
</details>
<br>
<details>
  <summary> Notifications screennshots </summary>

### Generic Notification example<br>
![Example 2.](/assets/example2.png)
![Example 3.](/assets/example3.png)
### OTP Notification example<br>
![Example 5.](/assets/example5.png)
</details>
<br>
🛡️ License
<br>GPL — use, modify, and deploy freely.