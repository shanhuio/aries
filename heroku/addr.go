package heroku

import (
	"os"
)

// Addr is the serving address for a Heroku web instance. Returns empty string
// if the PORT environment variable is missing.
func Addr() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ""
	}
	return ":" + port
}

// DBURL returns the database URL of the app.
func DBURL() string {
	return os.Getenv("DATABASE_URL")
}
