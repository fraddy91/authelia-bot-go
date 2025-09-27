package bot

import (
	"os"
	"strconv"
	"time"
)

var StartTime = time.Now()
var LastProcessed time.Time

var AllowRegistration = getEnvBool("ALLOW_REGISTRATION", true)

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}
