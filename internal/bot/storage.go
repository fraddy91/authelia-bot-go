package bot

import (
	"encoding/json"
	"os"
)

// ---------------------- Data Models ----------------------

type User struct {
	ChatID int64  `json:"chat_id"`
	Mode   string `json:"mode"`
	Admin  bool   `json:"admin"`
	Notify bool   `json:"notify"`
}

type PendingRequest struct {
	Email  string `json:"email"`
	Mode   string `json:"mode"`
	ChatID int64  `json:"chat_id"`
}

// ---------------------- File Paths ----------------------

var (
	usersFile   = "config/users.json"
	pendingFile = "config/pending.json"
	ignoreFile  = "config/ignore.json"
)

// ---------------------- Generic Helpers ----------------------

func loadJSON(path string, v interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	json.Unmarshal(data, v)
}

func saveJSON(path string, v interface{}) {
	os.MkdirAll("config", 0755)
	data, _ := json.MarshalIndent(v, "", "  ")
	os.WriteFile(path, data, 0644)
}

// ---------------------- Users ----------------------

func loadUsers() map[string]User {
	users := map[string]User{}
	loadJSON(usersFile, &users)
	return users
}

func saveUsers(users map[string]User) {
	saveJSON(usersFile, users)
}

// ---------------------- Pending ----------------------

func loadPending() map[string]PendingRequest {
	pending := map[string]PendingRequest{}
	loadJSON(pendingFile, &pending)
	return pending
}

func savePending(pending map[string]PendingRequest) {
	saveJSON(pendingFile, pending)
}

// ---------------------- Ignore ----------------------

func loadIgnore() map[string]bool {
	ignore := map[string]bool{}
	loadJSON(ignoreFile, &ignore)
	return ignore
}

func saveIgnore(ignore map[string]bool) {
	saveJSON(ignoreFile, ignore)
}

// ---------------------- Admin helpers ----------------------

func isAdmin(chatID int64) bool {
	users := loadUsers()
	for _, u := range users {
		if u.ChatID == chatID && u.Admin {
			return true
		}
	}
	return false
}
