package utils

import (
	"backend/config"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailType string

const (
	EmailDining EmailType = "dining"
	EmailEvents EmailType = "events"
	EmailPlay   EmailType = "play"
	EmailAdmin  EmailType = "admin"
)

type EmailSender struct {
	SMTPHost string
	SMTPPort int
	Email    string
	Password string
}

var senders map[EmailType]EmailSender

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func InitEmail(cfg *config.Config) {
	senders = map[EmailType]EmailSender{
		EmailDining: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    getEnvOrDefault("DINING_EMAIL", "dining@ticpin.in"),
			Password: getEnvOrDefault("DINING_APP_PASSWORD", ""),
		},
		EmailEvents: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    getEnvOrDefault("EVENTS_EMAIL", "events@ticpin.in"),
			Password: getEnvOrDefault("EVENTS_APP_PASSWORD", ""),
		},
		EmailPlay: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    getEnvOrDefault("PLAY_EMAIL", "play@ticpin.in"),
			Password: getEnvOrDefault("PLAY_APP_PASSWORD", ""),
		},
		EmailAdmin: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    getEnvOrDefault("ADMIN_EMAIL", "admin@ticpin.in"),
			Password: getEnvOrDefault("ADMIN_APP_PASSWORD", ""),
		},
	}
}

func SendEmail(emailType EmailType, to string, subject string, body string) error {
	if senders == nil {
		return fmt.Errorf("email system not initialized")
	}
	sender, ok := senders[emailType]
	if !ok {
		return fmt.Errorf("invalid email type: %s", emailType)
	}

	// Check if password is configured
	if sender.Password == "" {
		log.Printf("⚠️  Email sending failed: No password configured for %s (email: %s)", emailType, sender.Email)
		return fmt.Errorf("email account not configured: missing password for %s", emailType)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("TicPin <%s>", sender.Email))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(sender.SMTPHost, sender.SMTPPort, sender.Email, sender.Password)

	log.Printf("📧 Sending email to %s from %s...", to, sender.Email)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("❌ Failed to send email via %s to %s: %v", emailType, to, err)
		return err
	}

	log.Printf("✅ Email sent successfully to %s from %s", to, sender.Email)
	return nil
}

// SendOTPEmail sends an OTP email and falls back gracefully in development mode.
// In production it returns an error if sending fails (caller should return 500).
// In development it logs the OTP to the console so you can test without real SMTP.
func SendOTPEmail(emailType EmailType, to, subject, body, otp string) error {
	// In development: log the OTP to console so you can test even if SMTP is slow
	if config.CurrentEnv != "production" {
		log.Printf("💡  [DEV] OTP for %s → %s", to, otp)
	}

	err := SendEmail(emailType, to, subject, body)
	if err != nil {
		log.Printf("❌ [OTP] Failed to send OTP email to %s: %v", to, err)
		return err // always return the error so the controller can respond with 500
	}

	return nil
}

func UpdateSenderPassword(emailType EmailType, password string) {
	if sender, ok := senders[emailType]; ok {
		sender.Password = password
		senders[emailType] = sender
	}
}

func getEmailLogoSVG() string {
	return `<img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/logo.png?alt=media" width="140" height="35" style="display:block;margin:0 auto;">`
}

// getEmailSocialFooter returns the common email footer with help section, social icons and copyright
func getEmailSocialFooter() string {
	return `
<table width="611" align="center" style="margin-top:40px;color:#ffffff;">
<tr><td style="font-weight:600;font-size:35px;">LOOKING FOR HELP?</td></tr>
<tr><td style="padding-top:15px;font-size:20px;">Mail us at <span style="color:#4EA3FF;">support@ticpin.in</span> (10AM-5PM), and we’ll help you out.</td></tr>
<tr><td style="padding-top:30px;"><hr style="border:1px solid #ffffff;"></td></tr>
<tr>
<td align="center" style="padding:20px 0;">
<a href="#"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/whatsapp.png?alt=media" width="50" height="50"></a>
<a href="#" style="margin:0 25px;display:inline-block;"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/facebook.png?alt=media" width="50" height="50"></a>
<a href="#"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/insta.png?alt=media" width="50" height="50"></a>
<a href="#" style="margin:0 25px;display:inline-block;"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/youtube.png?alt=media" width="50" height="50"></a>
<a href="#"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/tiwteer.png?alt=media" width="50" height="50"></a>
</td>
</tr>
<tr><td><hr style="border:1px solid #ffffff;"></td></tr>
</table>`
}

func parseDateComponents(dateStr string) (day, date, month, dateMonth string) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Try dd Jan 2006
		t, err = time.Parse("02 Jan 2006", dateStr)
		if err != nil {
			return "", dateStr, "", dateStr
		}
	}
	return t.Format("Monday"), t.Format("02"), t.Format("January"), t.Format("02 January")
}

func GetOTPEmailTemplate(otp string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Ticpin OTP</title>
<link href="https://fonts.googleapis.com/css2?family=Anek+Latin:wght@500;600&display=swap" rel="stylesheet">
</head>

<body style="margin:0;padding:40px 20px;background:#f0f0f0;font-family:'Anek Latin', sans-serif;">

<table align="center" width="530" cellpadding="0" cellspacing="0" style="background:#0A0132;border-radius:15px;padding:20px 30px;">
<tr>
<td align="center">

<!-- White Card -->
<table width="450" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:15px;overflow:hidden;">

<!-- Purple Header -->
<tr>
<td align="center" style="margin-top: 0px; background:#5331EA;padding:25px 0;border-top-left-radius:15px;border-top-right-radius:15px;">
<img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/logo.png?alt=media&token=96a72021-6f7d-45c6-b270-77e73b80cf09" width="140" height="35" style="display:block;">
</td>
</tr>

<tr>
<td style="border-top:1px solid #AEAEAE;"></td>
</tr>

<!-- Content -->
<tr>
<td align="center" style="padding:35px;">

<div style="font-weight:600;font-size:30px;line-height:100%%;color:#000000;">
Welcome to Ticpin
</div>

<div style="margin-top:25px;font-weight:500;font-size:25px;line-height:28px;color:#000000;">
Your OTP for login is
</div>

<div style="margin-top:15px;font-weight:600;font-size:40px;line-height:50px;color:#000000;letter-spacing:5px;">
%s
</div>

<div style="margin-top:25px;font-weight:500;font-size:20px;line-height:22px;color:#000000;">
This is valid for 5mins
</div>

</td>
</tr>

</table>

<!-- Help Section -->
<table width="450" cellpadding="0" cellspacing="0" style="margin-top:20px;color:#ffffff;">
<tr>
<td style="font-weight:600;font-size:30px;line-height:100%%;">
LOOKING FOR HELP?
</td>
</tr>

<tr>
<td style="padding-top:10px;font-weight:500;font-size:20px;line-height:22px;">
Mail us at <span style="color:#4EA3FF;">support@ticpin.in</span> (10AM-5PM), and we’ll help you out.
</td>
</tr>

<tr>
<td style="padding-top:12px;">
<hr style="border:1px solid #ffffff;">
</td>
</tr>

<tr>
<td align="center" style="padding:12px 0;">

<a href="#">
    <img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/whatsapp.png?alt=media" width="50" height="50" style="display:inline-block;">
</a>
<a href="#" style="margin:0 25px;display:inline-block;">
    <img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/facebook.png?alt=media" width="50" height="50" style="display:inline-block;">
</a>
<a href="#">
    <img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/insta.png?alt=media" width="50" height="50" style="display:inline-block;">
</a>
<a href="#" style="margin:0 25px;display:inline-block;">
    <img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/youtube.png?alt=media" width="50" height="50" style="display:inline-block;">
</a>
<a href="#">
    <img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/tiwteer.png?alt=media" width="50" height="50" style="display:inline-block;">
</a>

</td>

</tr>

<tr>
<td>
<hr style="border:1px solid #ffffff;">
</td>
</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
	`, otp)
}

func GetPlayBookingEmailTemplate(playerName, venueName, sport, date, timeSlot, bookingID string) string {
	day, dayDate, month, dateMonth := parseDateComponents(date)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Play Booking Confirmed</title>
<link href="https://fonts.googleapis.com/css2?family=Anek+Latin:wght@500;600&display=swap" rel="stylesheet">
</head>
<body style="margin:0;padding:40px 20px;background:#f0f0f0;font-family:'Anek Latin', sans-serif;">
<table align="center" width="711" cellpadding="0" cellspacing="0" style="background:#0A0132;border-radius:15px;padding:40px 50px;">
<tr><td align="center">
<table width="611" cellpadding="0" cellspacing="0" style="background:#EBEBEB;border-radius:15px;overflow:hidden;">
<tr><td style="background:#5331EA;padding:35px;"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/logo.png?alt=media&token=96a72021-6f7d-45c6-b270-77e73b80cf09" width="165" height="41"></td></tr>
<tr><td style="border-top:1px solid #AEAEAE;"></td></tr>
<tr><td style="padding:30px 35px 10px 35px;">
<div style="font-weight:600;font-size:35px;color:#000;">Play booking confirmed <span style="color:#0AC655;">✔</span></div>
<div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">Booking Date : %s, %s, %s</div>
</td></tr>
<tr><td align="center" style="padding:20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:20px;">
<tr><td width="40%%"><div style="background:#AC9BF7;height:111px;width:197px;border-radius:8px;display:flex;align-items:center;justify-content:center;color:white;font-weight:600;">PLAY</div></td>
<td style="vertical-align:top;"><div style="font-weight:600;font-size:20px;color:#000;">%s</div><div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">%s</div></td></tr>
</table></td></tr>
<tr><td align="center" style="padding:0 35px 20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-size:17px;color:#686868;">Booking ID</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">#%s</td></tr>
<tr><td style="border-top:1px solid #D9D9D9;"></td></tr>
<tr><td style="font-size:17px;color:#686868;padding-top:15px;">Date & Time</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s %s %s | %s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Player Name</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Location</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Offer</td></tr>
<tr><td style="font-size:20px;color:#000;">None</td></tr>
</table></td></tr>
<tr><td align="center" style="padding:10px 35px 30px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-weight:600;font-size:20px;color:#000;padding-bottom:15px;">Notes</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Please arrive 10 minutes before your slot time.</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Carry a digital copy of this email for verification.</td></tr>
<tr><td style="font-size:15px;color:#686868;padding-top:15px;">See you there! <br> Team <span style="color:#5331EA;">Ticpin</span></td></tr>
</table></td></tr>
</table>
%s
</td></tr>
</table>
</body>
</html>`, day, dateMonth, timeSlot, venueName, sport, bookingID, day, dayDate, month, timeSlot, playerName, venueName, getEmailSocialFooter())
}

func GetDiningBookingEmailTemplate(guestName, restaurantName, date, timeSlot, bookingID string, guestCount int, specialRequest string) string {
	day, dayDate, month, dateMonth := parseDateComponents(date)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Table Booking Confirmed</title>
<link href="https://fonts.googleapis.com/css2?family=Anek+Latin:wght@500;600&display=swap" rel="stylesheet">
</head>
<body style="margin:0;padding:40px 20px;background:#f0f0f0;font-family:'Anek Latin', sans-serif;">
<table align="center" width="711" cellpadding="0" cellspacing="0" style="background:#0A0132;border-radius:15px;padding:40px 50px;">
<tr><td align="center">
<table width="611" cellpadding="0" cellspacing="0" style="background:#EBEBEB;border-radius:15px;overflow:hidden;">
<tr><td style="background:#5331EA;padding:35px;"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/logo.png?alt=media&token=96a72021-6f7d-45c6-b270-77e73b80cf09" width="165" height="41"></td></tr>
<tr><td style="border-top:1px solid #AEAEAE;"></td></tr>
<tr><td style="padding:30px 35px 10px 35px;">
<div style="font-weight:600;font-size:35px;color:#000;">Table booking confirmed <span style="color:#0AC655;">✔</span></div>
<div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">Booking Date : %s, %s, %s</div>
</td></tr>
<tr><td align="center" style="padding:20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:20px;">
<tr><td width="40%%"><div style="background:#AC9BF7;height:111px;width:197px;border-radius:8px;display:flex;align-items:center;justify-content:center;color:white;font-weight:600;">DINING</div></td>
<td style="vertical-align:top;"><div style="font-weight:600;font-size:20px;color:#000;">%s</div><div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">%s</div></td></tr>
</table></td></tr>
<tr><td align="center" style="padding:0 35px 20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-size:17px;color:#686868;">Booking ID</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">#%s</td></tr>
<tr><td style="border-top:1px solid #D9D9D9;"></td></tr>
<tr><td style="font-size:17px;color:#686868;padding-top:15px;">Date & Time</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s %s %s %s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Number of guest(s)</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%d</td></tr>
<tr><td style="font-size:17px;color:#686868;">Location</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Offer</td></tr>
<tr><td style="font-size:20px;color:#000;">None</td></tr>
</table></td></tr>
<tr><td align="center" style="padding:10px 35px 30px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-weight:600;font-size:20px;color:#000;padding-bottom:15px;">Notes</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Please arrive 10 minutes before your reserved time.</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Late arrivals may cause reservation cancellation.</td></tr>
<tr><td style="font-size:15px;color:#686868;padding-top:15px;">See you there! <br> Team <span style="color:#5331EA;">Ticpin</span></td></tr>
</table></td></tr>
</table>
%s
</td></tr>
</table>
</body>
</html>`, day, dateMonth, timeSlot, restaurantName, restaurantName, bookingID, day, dayDate, month, timeSlot, guestCount, restaurantName, getEmailSocialFooter())
}

func GetEventBookingEmailTemplate(guestName, eventName, venue, date, timeStr string, ticketCount int, qrImageURL string, bookingID string) string {
	day, dayDate, month, dateMonth := parseDateComponents(date)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Event Booking Confirmed</title>
<link href="https://fonts.googleapis.com/css2?family=Anek+Latin:wght@500;600&display=swap" rel="stylesheet">
</head>
<body style="margin:0;padding:40px 20px;background:#f0f0f0;font-family:'Anek Latin', sans-serif;">
<table align="center" width="711" cellpadding="0" cellspacing="0" style="background:#0A0132;border-radius:15px;padding:40px 50px;">
<tr><td align="center">
<table width="611" cellpadding="0" cellspacing="0" style="background:#EBEBEB;border-radius:15px;overflow:hidden;">
<tr><td style="background:#5331EA;padding:35px;"><img src="https://firebasestorage.googleapis.com/v0/b/ticpin-fa6d2.firebasestorage.app/o/logo.png?alt=media&token=96a72021-6f7d-45c6-b270-77e73b80cf09" width="165" height="41"></td></tr>
<tr><td style="border-top:1px solid #AEAEAE;"></td></tr>
<tr><td style="padding:30px 35px 10px 35px;">
<div style="font-weight:600;font-size:35px;color:#000;">Event booking confirmed <span style="color:#0AC655;">✔</span></div>
<div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">Booking Date : %s, %s, %s</div>
</td></tr>
<tr><td align="center" style="padding:20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:20px;">
<tr><td width="40%%"><div style="background:#AC9BF7;height:111px;width:197px;border-radius:8px;display:flex;align-items:center;justify-content:center;color:white;font-weight:600;">EVENT</div></td>
<td style="vertical-align:top;"><div style="font-weight:600;font-size:20px;color:#000;">%s</div><div style="margin-top:10px;font-weight:500;font-size:20px;color:#686868;">%s</div></td></tr>
</table></td></tr>
<tr><td align="center" style="padding:0 35px 20px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-size:17px;color:#686868;">Booking ID</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">#%s</td></tr>
<tr><td style="border-top:1px solid #D9D9D9;"></td></tr>
<tr><td style="font-size:17px;color:#686868;padding-top:15px;">Date & Time</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s %s %s | %s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Number of ticket(s)</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%d</td></tr>
<tr><td style="font-size:17px;color:#686868;">Location</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Gate opening time</td></tr>
<tr><td style="font-size:20px;color:#000;padding-bottom:15px;">%s</td></tr>
<tr><td style="font-size:17px;color:#686868;">Offer</td></tr>
<tr><td style="font-size:20px;color:#000;">None</td></tr>
</table></td></tr>
<tr><td align="center" style="padding:30px;">
<div style="width:215px;height:215px;background:rgba(83,49,234,0.15);display:flex;align-items:center;justify-content:center;border-radius:12px;">
<img src="%s" width="180">
</div>
</td></tr>
<tr><td align="center" style="padding:10px 35px 30px 35px;">
<table width="100%%" style="background:#FFFFFF;border-radius:10px;padding:25px;">
<tr><td style="font-weight:600;font-size:20px;color:#000;padding-bottom:15px;">Notes</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Please arrive 15 minutes before the event start time.</td></tr>
<tr><td style="font-size:17px;color:#686868;padding-bottom:10px;">◆ Follow all venue rules and safety instructions.</td></tr>
<tr><td style="font-size:15px;color:#686868;padding-top:15px;">See you there! <br> Team <span style="color:#5331EA;">Ticpin</span></td></tr>
</table></td></tr>
</table>
%s
</td></tr>
</table>
</body>
</html>`, day, dateMonth, timeStr, eventName, venue, bookingID, day, dayDate, month, timeStr, ticketCount, venue, timeStr, qrImageURL, getEmailSocialFooter())
}

// GetPassPurchaseEmailTemplate returns the pass purchase confirmation email HTML
func GetPassPurchaseEmailTemplate(name, passID, purchaseDate, expiryDate string, amount int) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Ticpin Pass Purchase Confirmation</title>
<link href="https://fonts.googleapis.com/css2?family=Anek+Latin:wght@500;600&display=swap" rel="stylesheet">
</head>

<body style="margin:0;padding:40px 20px;background:#f0f0f0;font-family:'Anek Latin', sans-serif;">

<table align="center" width="530" cellpadding="0" cellspacing="0" style="background:#0A0132;border-radius:15px;padding:20px 30px;">
<tr>
<td align="center">

<!-- White Card -->
<table width="450" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:15px;overflow:hidden;">

<!-- Purple Header -->
<tr>
<td align="center" style="background:linear-gradient(135deg, #5331EA 0%%, #E91E63 100%%);padding:40px 0;border-top-left-radius:15px;border-top-right-radius:15px;">
%s
<div style="margin-top:15px;font-weight:600;font-size:28px;line-height:100%%;color:#ffffff;">
🎉 Welcome to Ticpin Pass!
</div>
<div style="margin-top:8px;font-weight:500;font-size:18px;line-height:24px;color:#ffffff;">
Your Premium Membership is Now Active
</div>
</td>
</tr>

<!-- Content -->
<tr>
<td style="padding:35px;">

<div style="font-weight:600;font-size:22px;line-height:100%%;color:#000000;">
Hey %s! 👋
</div>

<div style="margin-top:12px;font-weight:500;font-size:15px;line-height:22px;color:#686868;">
Your Ticpin Pass purchase was successful. Get ready to enjoy exclusive benefits on your favorite activities.
</div>

<!-- Pass Details Card -->
<table width="100%%" cellpadding="0" cellspacing="0" style="margin:20px 0;background:#f8f4ff;border-radius:10px;padding:20px;border:1px solid #e8dff5;">
<tr>
<td style="font-weight:600;font-size:16px;color:#000;padding-bottom:15px;">Pass Details</td>
</tr>
<tr>
<td>
<table width="100%%" cellpadding="0" cellspacing="0">
<tr style="border-bottom:1px solid #eee;">
<td style="padding:10px 0;font-size:14px;color:#686868;">Pass ID</td>
<td style="padding:10px 0;font-weight:600;color:#000;text-align:right;">%s</td>
</tr>
<tr style="border-bottom:1px solid #eee;">
<td style="padding:10px 0;font-size:14px;color:#686868;">Amount Paid</td>
<td style="padding:10px 0;font-weight:600;color:#5331EA;text-align:right;">₹%d</td>
</tr>
<tr style="border-bottom:1px solid #eee;">
<td style="padding:10px 0;font-size:14px;color:#686868;">Purchase Date</td>
<td style="padding:10px 0;font-weight:500;color:#000;text-align:right;">%s</td>
</tr>
<tr>
<td style="padding:10px 0;font-size:14px;color:#686868;">Valid Until</td>
<td style="padding:10px 0;font-weight:600;color:#E91E63;text-align:right;">%s</td>
</tr>
</table>
</td>
</tr>
</table>

<!-- Benefits Section -->
<div style="margin-top:20px;">
<div style="font-weight:600;font-size:16px;color:#000;margin-bottom:12px;">Your Benefits</div>
<table width="100%%" cellpadding="0" cellspacing="0">
<tr>
<td style="width:33%%;padding:10px;vertical-align:top;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#fff;border-radius:8px;padding:12px;border-left:4px solid #5331EA;text-align:center;">
<tr><td style="font-size:20px;">🎾</td></tr>
<tr><td style="font-weight:600;font-size:13px;color:#000;margin-top:5px;">2 Free Turf Bookings</td></tr>
<tr><td style="font-size:11px;color:#999;margin-top:3px;">Valid for 3 months</td></tr>
</table>
</td>
<td style="width:33%%;padding:10px;vertical-align:top;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#fff;border-radius:8px;padding:12px;border-left:4px solid #5331EA;text-align:center;">
<tr><td style="font-size:20px;">⚡</td></tr>
<tr><td style="font-weight:600;font-size:13px;color:#000;margin-top:5px;">15%% Discount</td></tr>
<tr><td style="font-size:11px;color:#999;margin-top:3px;">On all bookings</td></tr>
</table>
</td>
<td style="width:33%%;padding:10px;vertical-align:top;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#fff;border-radius:8px;padding:12px;border-left:4px solid #5331EA;text-align:center;">
<tr><td style="font-size:20px;">🎯</td></tr>
<tr><td style="font-weight:600;font-size:13px;color:#000;margin-top:5px;">Priority Support</td></tr>
<tr><td style="font-size:11px;color:#999;margin-top:3px;">Dedicated help</td></tr>
</table>
</td>
</tr>
</table>
</div>

<!-- Reminder -->
<div style="margin-top:20px;background:#fff3cd;padding:12px;border-radius:5px;border-left:4px solid #FF9800;">
<div style="font-size:13px;color:#856404;">
⚠️ <strong>Reminder:</strong> Your pass expires on %s. Renew before expiry to avoid losing your benefits!
</div>
</div>

</td>
</tr>

<!-- Footer -->
<tr>
<td align="center" style="padding:20px;border-top:1px solid #eee;">
<div style="font-size:13px;color:#999;margin-bottom:8px;">
Questions? <a href="mailto:support@ticpin.com" style="color:#5331EA;text-decoration:none;font-weight:600;">Contact our support team</a>
</div>
<div style="font-size:12px;color:#999;">
© 2026 Ticpin. All rights reserved.
</div>
</td>
</tr>

</table>

</td>
</tr>
</table>
</body>
</html>
	`, getEmailLogoSVG(), name, passID, amount, purchaseDate, expiryDate, expiryDate)
}
