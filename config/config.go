package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	FirebaseCredentials string
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
	// Payment gateway credentials
	CashfreePaymentAppID  string
	CashfreePaymentSecret string
	RazorpayKeyID         string
	RazorpayKeySecret     string
	GroqAPIKey            string
	FirebaseKeyPath       string
}

// CurrentEnv is set at startup and read by other packages to check the running environment.
var CurrentEnv = "development"

func LoadConfig() *Config {
	// Try loading local .env first
	if err := godotenv.Load(); err != nil {
		// Fallback to Render's secret path
		if err := godotenv.Load("/etc/secrets/.env"); err != nil {
			log.Println("No .env file found in root or /etc/secrets/, using system environment variables")
		} else {
			log.Println("Loaded environment variables from /etc/secrets/.env")
		}
	}

	cfg := &Config{
		Port:                  getEnv("PORT", ":9000"),
		Env:                   getEnv("ENV", "development"),
		FirebaseCredentials:   getEnv("FIREBASE_CREDENTIALS", ""),
		DiningEmail:           getEnv("DINING_EMAIL", "dining@ticpin.in"),
		DiningAppPass:         getEnv("DINING_APP_PASSWORD", ""),
		EventsEmail:           getEnv("EVENTS_EMAIL", "events@ticpin.in"),
		EventsAppPass:         getEnv("EVENTS_APP_PASSWORD", ""),
		PlayEmail:             getEnv("PLAY_EMAIL", "play@ticpin.in"),
		PlayAppPass:           getEnv("PLAY_APP_PASSWORD", ""),
		AdminEmail:            getEnv("ADMIN_EMAIL", "admin@ticpin.in"),
		AdminAppPass:          getEnv("ADMIN_APP_PASSWORD", ""),
		CashfreeClientID:      getEnv("CASHFREE_CLIENT_ID", ""),
		CashfreeSecret:        getEnv("CASHFREE_CLIENT_SECRET", ""),
		CashfreePANURL:        getEnv("CASHFREE_PAN_VERIFY_URL", "https://api.cashfree.com/verification/pan-lite"),
		CashfreePANGSTINURL:   getEnv("CASHFREE_PAN_GSTIN_URL", "https://api.cashfree.com/verification/pan-gstin"),
		CashfreePaymentAppID:  getEnv("CASHFREE_PAYMENT_APP_ID", ""),
		CashfreePaymentSecret: getEnv("CASHFREE_PAYMENT_SECRET", ""),
		RazorpayKeyID:         getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret:     getEnv("RAZORPAY_KEY_SECRET", ""),
		GroqAPIKey:            getEnv("GROQ_API_KEY", ""),
		FirebaseKeyPath:       getEnv("FIREBASE_KEY_PATH", "ticpin-fa6d2-firebase-adminsdk-fbsvc-53d16fed36.json"),
	}

	CurrentEnv = cfg.Env
	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
