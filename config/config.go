package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	FirebaseCredentials string // Full service account JSON
	DiningEmail         string
	DiningAppPass       string
	EventsEmail         string
	EventsAppPass       string
	PlayEmail           string
	PlayAppPass         string
	AdminEmail          string
	AdminAppPass        string
	CashfreeClientID    string
	CashfreeSecret      string
	CashfreePANURL      string
	CashfreePANGSTINURL string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Port:                getEnv("PORT", ":9000"),
		Env:                 getEnv("ENV", "development"),
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),
		DiningEmail:         getEnv("DINING_EMAIL", "dining@ticpin.in"),
		DiningAppPass:       getEnv("DINING_APP_PASSWORD", ""),
		EventsEmail:         getEnv("EVENTS_EMAIL", "events@ticpin.in"),
		EventsAppPass:       getEnv("EVENTS_APP_PASSWORD", ""),
		PlayEmail:           getEnv("PLAY_EMAIL", "play@ticpin.in"),
		PlayAppPass:         getEnv("PLAY_APP_PASSWORD", "irjs ojvn fbmh orht"),
		AdminEmail:          getEnv("ADMIN_EMAIL", "admin@ticpin.in"),
		AdminAppPass:        getEnv("ADMIN_APP_PASSWORD", ""),
		CashfreeClientID:    getEnv("CASHFREE_CLIENT_ID", ""),
		CashfreeSecret:      getEnv("CASHFREE_CLIENT_SECRET", ""),
		CashfreePANURL:      getEnv("CASHFREE_PAN_VERIFY_URL", "https://api.cashfree.com/verification/pan/advance"),
		CashfreePANGSTINURL: getEnv("CASHFREE_PAN_GSTIN_URL", "https://api.cashfree.com/verification/pan-gstin"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
