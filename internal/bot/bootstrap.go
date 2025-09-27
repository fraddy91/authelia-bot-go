package bot

import (
	"log"
	"os"
)

func EnsureStartupFiles() {
	// Ensure folders exist
	os.MkdirAll("config", 0755)
	os.MkdirAll("notifications", 0755)

	// Ensure files exist
	files := map[string][]byte{
		"config/users.json":              []byte("{}"),
		"config/pending.json":            []byte("{}"),
		"config/ignore.json":             []byte("{}"),
		"notifications/notification.txt": []byte(""),
	}

	for path, content := range files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.WriteFile(path, content, 0644)
			if err != nil {
				log.Printf("⚠️ Failed to create %s: %v", path, err)
			} else {
				log.Printf("📁 Created missing file: %s", path)
			}
		}
	}
}
