package utils

import (
	"fmt"
	"log"
)

// SendSMS is a placeholder for sending SMS
// You can integrate services like Twilio, MessageBird, or Firebase here.
func SendSMS(phone string, message string) error {
	// For now, we just log the OTP to the console.
	// In production, connect to a real SMS gateway.
	log.Printf("\n--- SMS TO %s ---\n%s\n-------------------\n", phone, message)

	// Example Twilio implementation (commented out):
	/*
		accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
		authToken := os.Getenv("TWILIO_AUTH_TOKEN")
		from := os.Getenv("TWILIO_PHONE_NUMBER")

		url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

		data := url.Values{}
		data.Set("To", phone)
		data.Set("From", from)
		data.Set("Body", message)

		req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		return err
	*/

	return nil
}

func GetOTPSMSMessage(otp string) string {
	return fmt.Sprintf("Your TicPin verification code is: %s. Valid for 5 minutes.", otp)
}
