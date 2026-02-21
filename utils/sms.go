package utils

import (
	"fmt"
	"log"
)

func SendSMS(phone string, message string) error {
	
	log.Printf("\n--- SMS TO %s ---\n%s\n-------------------\n", phone, message)

	return nil
}

func GetOTPSMSMessage(otp string) string {
	return fmt.Sprintf("Your TicPin verification code is: %s. Valid for 5 minutes.", otp)
}
