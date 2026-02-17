package utils

import (
	"backend/config"
	"crypto/tls"
	"fmt"
	"log"

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

func InitEmail(cfg *config.Config) {
	senders = map[EmailType]EmailSender{
		EmailDining: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    cfg.DiningEmail,
			Password: cfg.DiningAppPass,
		},
		EmailEvents: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    cfg.EventsEmail,
			Password: cfg.EventsAppPass,
		},
		EmailPlay: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    cfg.PlayEmail,
			Password: cfg.PlayAppPass,
		},
		EmailAdmin: {
			SMTPHost: "smtp.gmail.com",
			SMTPPort: 587,
			Email:    cfg.AdminEmail,
			Password: cfg.AdminAppPass,
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

	m := gomail.NewMessage()
	m.SetHeader("From", sender.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(sender.SMTPHost, sender.SMTPPort, sender.Email, sender.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email via %s: %v", emailType, err)
		return err
	}

	return nil
}

// UpdateSenderPassword allows updating passwords dynamically if needed
func UpdateSenderPassword(emailType EmailType, password string) {
	if sender, ok := senders[emailType]; ok {
		sender.Password = password
		senders[emailType] = sender
	}
}

func GetOTPEmailTemplate(otp string) string {
	return fmt.Sprintf(`
<div style="margin: 0; padding: 0; font-family: 'Anek Latin', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">
    <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #0A0132; padding: 40px 10px;">
        <tr>
            <td align="center">
                <!-- Main Container -->
                <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 480px; margin: 0 auto; background-color: #0A0132; border-radius: 24px; overflow: hidden; border: 1px solid rgba(255,255,255,0.1);">
                    <!-- Header -->
                    <tr>
                        <td align="center" style="background-color: #5331EA; height: 100px; color: #ffffff;">
                            <h1 style="margin: 0; font-size: 28px; font-weight: 700; letter-spacing: 2px;">TICPIN</h1>
                        </td>
                    </tr>
                    
                    <!-- Content Card (White) -->
                    <tr>
                        <td align="center" style="padding: 24px;">
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 420px; background-color: #ffffff; border-radius: 20px; margin: 0 auto;">
                                <tr>
                                    <td align="center" style="padding: 40px 20px;">
                                        <h2 style="color: #000000; font-size: 28px; font-weight: 600; margin: 0 0 32px 0; line-height: 32px;">Welcome to Ticpin</h2>
                                        
                                        <p style="color: #000000; font-size: 18px; font-weight: 500; margin: 0 0 12px 0; line-height: 22px;">Your OTP for login is</p>
                                        
                                        <div style="color: #5331EA; font-size: 40px; font-weight: 700; margin: 0 0 24px 0; line-height: 48px; letter-spacing: 4px;">%s</div>
                                        
                                        <p style="color: #666666; font-size: 14px; font-weight: 500; margin: 0; line-height: 18px;">This is valid for 5 mins</p>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <!-- Footer Help Section -->
                    <tr>
                        <td style="padding: 32px 32px 48px 32px;">
                            <h3 style="color: #ffffff; font-size: 35px; font-weight: 600; margin: 0 0 24px 0; line-height: 38px; text-transform: uppercase;">LOOKING FOR HELP?</h3>
                            
                            <p style="color: #ffffff; font-size: 20px; font-weight: 500; margin: 0 0 32px 0; line-height: 22px;">
                                Mail us at <a href="mailto:support@ticpin.in" style="color: #60a5fa; text-decoration: underline;">support@ticpin.in</a> (10AM-5PM), and we'll help you out.
                            </p>
                            
                            <!-- Divider -->
                            <div style="border-top: 1px solid rgba(255,255,255,0.2); margin-bottom: 32px;"></div>

                            <!-- Social Icons -->
                            <table border="0" cellpadding="0" cellspacing="0" align="center">
                                <tr>
                                    <td style="padding: 0 12px;">
                                        <a href="#"><img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNTAiIGhlaWdodD0iNTAiIHZpZXdCb3g9IjAgMCA1MCA1MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik0zNC41ODI3IDI5LjE2NjZDNDM0LjE2NiAyOC45NTgzIDMxLjQ1NzcgMjcuNzA4MyAzMS4wNDEgMjcuNUMzMC42MjQzIDI3LjI5MTYgMzAuMjA3NyAyNy4yOTE2IDI5Ljc5MSAyNy43MDgzQzI5LjM3NDMgMjguMTI1IDI4LjU0MSAyOS4zNzUgMjguMTI0MyAyOS43OTE2QzI3LjkxNiAzMC4yMDgzIDI3LjQ5OTMgMzAuMjA4MyAyNy4wODI3IDMwQzI1LjYyNDMgMjkuMzc1IDI0LjE2NiAyOC41NDE2IDIyLjkxNiAyNy41QzIxLjg3NDMgMjYuNDU4MyAyMC44MzI3IDI1LjIwODMgMTkuOTk5MyAyMy45NTgzQzE5Ljc5MSAyMy41NDE2IDE5Ljk5OTMgMjMuMTI1IDIwLjIwNzcgMjIuOTE2NkMyMC40MTYgMjIuNzA4MyAyMC42MjQzIDIyLjI5MTYgMjEuMDQxIDIyLjA4MzNDMjEuMjQ5MyAyMS44NzUgMjEuNDU3NyAyMS40NTgzIDIxLjQ1NzcgMjEuMjVDMjEuNjY2IDIxLjA0MTYgMjEuNjY2IDIwLjYyNSAyMS40NTc3IDIwLjQxNjZDMjEuMjQ5MyAyMC4yMDgzIDIwLjIwNzcgMTcuNzA4MyAxOS43OTEgMTYuNjY2NkMxOS41ODI3IDE1LjIwODMgMTkuMTY2IDE1LjIwODMgMTguNzQ5MyAxNS4yMDgzSDE3LjcwNzdDMTcuMjkxIDE1LjIwODMgMTYuNjY2IDE1LjYyNSAxNi40NTc3IDE1LjgzMzNDMTUuMjA3NyAxNy4wODMzIDE0LjU4MjcgMTguNTQxNiAxNC41ODI3IDIwLjIwODNDMTQuNzkxIDIyLjA4MzMgMTUuNDE2IDIzLjk1ODMgMTYuNjY2IDI1LjYyNUMxOC45NTc3IDI4Ljk1ODMgMjEuODc0MyAzMS42NjY2IDI1LjQxNiAzMy4zMzMzQzI2LjQ1NzcgMzMuNzUgMjcuMjkxIDM0LjE2NjYgMjguMzMyNyAzNC4zNzVDMjkuMzc0MyAzNC43OTE2IDMwLjQxNiAzNC43OTE2IDMxLjY2NiAzNC41ODMzQzMzLjEyNDMgMzQuMzc1IDM0LjM3NDMgMzMuMzMzMyAzNS4yMDc3IDMyLjA4MzNDMzUuNjI0MyAzMS4yNSAzNS42MjQzIDMwLjQxNjYgMzUuNDE2IDI5LjU4MzNMMzQuNTgyNyAyOS4xNjY2Wk0zOS43OTEgMTAuMjA4M0MzMS42NjYgMi4wODMzMSAxOC41NDEgMi4wODMzMSAxMC40MTYgMTAuMjA4M0MzLjc0OTM1IDE2Ljg3NSAyLjQ5OTM1IDI3LjA4MzMgNy4wODI2OCAzNS4yMDgzTDQuMTY2MDIgNDUuODMzM0wxNS4yMDc3IDQyLjkxNjZDMTguMzMyNyA0NC41ODMzIDIxLjY2NiA0NS40MTY2IDI0Ljk5OTMgNDUuNDE2NkMzNi40NTc3IDQ1LjQxNjYgNDUuNjI0MyAzNi4yNSA0NS42MjQzIDI0Ljc5MTZDNDUuODMyNyAxOS4zNzUgNDMuNTQxIDE0LjE2NjYgMzkuNzkxIDEwLjIwODNaTTM0LjE2NiAzOS4zNzVDMzEuNDU3NyA0MS4wNDE2IDI4LjMzMjcgNDIuMDgzMyAyNC45OTkzIDQyLjA4MzNDMjEuODc0MyA0Mi4wODMzIDE4Ljk1NzcgNDEuMjUgMTYuMjQ5MyAzOS43OTE2TDE1LjYyNDMgMzkuMzc1TDkuMTY2MDIgNDEuMDQxNkwxMC44MzI3IDM0Ljc5MTZMMTAuNDE2IDM0LjE2NjZDNS40MTYwMiAyNS44MzMzIDcuOTE2MDEgMTUuNDE2NiAxNi4wNDEgMTAuMjA4M0MyNC4xNjYgNC45OTk5OCAzNC41ODI3IDcuNzA4MzEgMzkuNTgyNyAxNS42MjVDNDQuNTgyNyAyMy43NSA0Mi4yOTEgMzQuMzc1IDM0LjE2NiAzOS4zNzVaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4=" width="40" height="40"></a>
                                    </td>
                                    <td style="padding: 0 12px;">
                                        <a href="#"><img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNTAiIGhlaWdodD0iNTAiIHZpZXdCb3g9IjAgMCA1MCA1MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik0yNSA1QzEzLjk1NSA1IDUgMTMuOTU1IDUgMjVDNSAzNS4wMjY3IDEyLjM4NjcgNDMuMzA2NyAyMi4wMSA0NC43NTMzVjMwLjNIMTcuMDYxN1YyNS4wNDMzSDIyLjAxVjIxLjU0NUMyMi4wMSAxNS43NTMzIDI0LjgzMTcgMTMuMjExNyAyOS42NDUgMTMuMjExN0MzMS45NSAxMy4yMTE3IDMzLjE3IDEzLjM4MzMgMzMuNzQ2NyAxMy40NlYxOC4wNDgzSDMwLjQ2MzNDMjguNDIgMTguMDQ4MyAyNy43MDY3IDE5Ljk4NjcgMjcuNzA2NyAyMi4xN1YyNS4wNDMzSDMzLjY5NUwzMi44ODMzIDMwLjNIMjcuNzA2N1Y0NC43OTVDMzcuNDY4MyA0My40NzE3IDQ1IDM1LjEyNSA0NSAyNUM0NSAxMy45NTUgMzYuMDQ1IDUgMjUgNVoiIGZpbGw9IndoaXRlIi8+Cjwvc3ZnPg==" width="40" height="40"></a>
                                    </td>
                                    <td style="padding: 0 12px;">
                                        <a href="#"><img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNTAiIGhlaWdodD0iNTAiIHZpZXdCb3g9IjAgMCA1MCA1MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik0zNi4xMjQ0IDExLjM3NUMzNS42Mjk5IDExLjM3NSAzNS4xNDY1IDExLjUyMTYgMzQuNzM1NCAxMS43OTYzQzM0LjMyNDMgMTIuMDcxMSAzNC4wMDM5IDEyLjQ2MTUgMzMuODE0NyAxMi45MTgzQzMzLjYyNTQgMTMuMzc1MSAzMy41NzU5IDEzLjg3NzggMzMuNjcyNCAxNC4zNjI3QzMzLjc2ODkgMTQuODY3NyAzNC4wMDcgMTUuMjkzMiAzNC4zNTY2IDE1LjY0MjhDMzQuNzA2MiAxNS45OTI0IDM1LjE1MTcgMTYuMjMwNSAzNS42MzY2IDE2LjMyNzBDMzYuMTIxNiAxNi40MjM0IDM2LjYyNDIgMTYuMzczOSAzNy4wODExIDE2LjE4NDdDMzcuNTM3OSAxNS45OTU1IDM3LjkyODMgMTUuNjc1MSAzOC4yMDMgMTUuMjYzOUMzOC40Nzc3IDE0Ljg1MjggMzguNjI0NCAxNC4zNjk1IDM4LjYyNDQgMTMuODc1QzM4LjYyNDQgMTMuMjEyIDM4LjM2MSAxMi41NzYxIDM3Ljg5MjEgMTIuMTA3M0MzNy40MjMzIDExLjYzODQgMzYuNzg3NCAxMS4zNzUgMzYuMTI0NCAxMS4zNzVNNDUuNzA3NyAxNi40MTY3QzQ1LjY2NzIgMTQuNjg4MSA0NS4zNDM0IDEyLjk3NzkgNDQuNzQ5NCAxMS4zNTQyQzQ0LjIxOTYgOS45NjQ4OCA0My4zOTUyIDguNzA2ODUgNDIuMzMyNyA3LjY2NjY5QzQxLjMwMTEgNi41OTg4MiA0MC4wNDAxIDUuNzc5NTUgMzguNjQ1MiA1LjI3MDg1QzM3LjAyNTcgNC42NTg2OSAzNS4zMTM2IDQuMzI3NTQgMzMuNTgyNyA0LjI5MTY5QzMxLjM3NDQgNC4xNjY2OSAzMC42NjYgNC4xNjY2OSAyNC45OTk0IDQuMTY2NjlDMTkuMzMyNyA0LjE2NjY5IDE4LjYyNDQgNC4xNjY2OSAxNi40MTYgNC4yOTE2OUMxNC42ODUxIDQuMzI3NTQgMTIuOTczIDQuNjU4NjkgMTEuMzUzNSA1LjI3MDg1QzkuOTYxMTcgNS43ODQ3IDguNzAxMjkgNi42MDMyNyA3LjY2NjAyIDcuNjY2NjlDNi41OTgxNSA4LjY5ODMxIDUuNzc4ODggOS45NTkyNyA1LjI3MDE4IDExLjM1NDJDNC42NTgwMiAxMi45NzM3IDQuMzI2ODcgMTQuNjg1OCA0LjI5MTAyIDE2LjQxNjdDNC4xNjYwMiAxOC42MjUgNC4xNjYwMiAxOS4zMzM0IDQuMTY2MDIgMjVDNC4xNjYwMiAzMC42NjY3IDQuMTY2MDIgMzEuMzc1IDQuMjkxMDIgMzMuNTgzNEM0LjMyNjg3IDM1LjMxNDMgNC42NTgwMiAzNy4wMjY0IDUuMjcwMTggMzguNjQ1OUM1Ljc3ODg4IDQwLjA0MDggNi41OTgxNSA0MS4zMDE3IDcuNjY2MDIgNDIuMzMzNEM4LjcwMTI5IDQzLjM5NjggOS45NjExNyA0NC4yMTUzIDExLjM1MzUgNDQuNzI5MkMxMi45NzMwIDQ1LjM0MTQgMTQuNjg1MSA0NS42NzI1IDE2LjQxNiA0NS43MDg0QzE4LjYyNDQgNDUuODMzNCAxOS4zMzI3IDQ1LjgzMzQgMjQuOTk5NCA0NS44MzM0QzMwLjY2NiA0NS44MzM0IDMxLjM3NDQgNDUuODMzNCAzMy41ODI3IDQ1LjcwODRDMzUuMzEzNiA0NS42NzI1IDM3LjAyNTcgNDUuMzQxNCAzOC42NDUyIDQ0LjcyOTJDNDAuMDQwMSA0NC4yMjA1IDQxLjMwMTEgNDMuNDAxMiA0Mi4zMzI3IDQyLjMzMzRDNDMuMzk5OSA0MS4yOTcxIDQ0LjIyNTEgNDAuMDM3OSA0NC43NDk0IDM4LjY0NTlDNDUuMzQzNCAzNy4wMjIxIDQ1LjY2NzIgMzUuMzExOSA0NS43MDc3IDMzLjU4MzRDNDUuNzA3NyAzMS4zNzUgNDUuODMyNyAzMC42NjY3IDQ1LjgzMjcgMjVDNDUuODMyNyAxOS4zMzM0IDQ1LjgzMjcgMTguNjI1IDQ1LjcwNzcgMTYuNDE2N000MS45NTc3IDMzLjMzMzRDNDEuOTQyNSAzNC42NTU4IDQxLjcwMyAzNS45NjYxIDQxLjI0OTQgMzcuMjA4NEM0MC45MTY3IDM4LjExNTAgNDAuMzgyNCAzOC45MzQyIDM5LjY4NjkgMzkuNjA0MkMzOS4wMTExIDQwLjI5MjcgMzguMTkzNiA0MC44MjU5IDM3LjI5MTcgNDEuMTY2N0MzNi4wNDg3IDQxLjYyMDMgMzQuNzM4NSA0MS44NTk5IDMzLjQxNiA0MS44NzVDMzEuMzMyNyA0MS45NzkyIDMwLjU2MTkgNDIgMjUuMDgyNyA0MkMxOS42MDM1IDQyIDE4LjgzMjcgNDIgMTYuNzQ5MyA0MS44NzVDMTUuNDI2OCA0MS44NTk5IDE0LjExNjYgNDEuNjIwMyAxMi44NzQzIDQxLjE2NjdDMTEuOTcwNyA0MC44MjY3IDExLjE1MTEgNDAuMjkzMyAxMC40NzYgMzkuNjA0MkM5Ljc4MDIzIDM4LjkzNDIgOS4yNDU5OCAzOC4xMTUwIDguOTE2MDIgMzcuMjA4NEM4LjQ2MjQ0IDM1Ljk2NjEgOC4yMjI4NiAzNC42NTU4IDguMjA3NjkgMzMuMzMzNEM4LjEwMzUyIDMxLjI1IDguMTAzNTIgMzAuNSA4LjEwMzUyIDI1QzguMTAzNTIgMTkuNSA4LjEwMzUyIDE4Ljc1IDguMjA3NjkgMTYuNjY2N0M4LjIyMjg2IDE1LjM0NDIgOC40NjI0NCAxNC4wMzQgOC45MTYwMiAxMi43OTE3QzkuMjQ1OTggMTEuODg1IDkuNzgwMjMgMTEuMDY1OCAxMC40NzYgMTAuMzk1OUMxMS4xNTExIDkuNzA2NzcgMTEuOTcwNyA5LjE3MzM5IDEyLjg3NDMgOC44MzMzNkMxNC4xMTY2IDguMzc5NzcgMTUuNDI2OCA4LjE0MDE5IDE2Ljc0OTMgOC4xMjUwMkMxOC44MzI3IDguMDIwODYgMTkuNjAzNSA4LjAyMDg2IDI1LjA4MjcgOC4wMjA4NkMzMC41NjE5IDguMDIwODYgMzEuMzMyNyA4LjAyMDg2IDMzLjQxNiA4LjEyNTAyQzM0LjczODUgOC4xNDAxOSAzNi4wNDg3IDguMzc5NzcgMzcuMjkxIDguODMzMzZDMzguMTk0NiA5LjE3MzM5IDM5LjAxNDIgOS43MDY3NyAzOS42ODkzIDEwLjM5NTlDNDAuMzg1MSAxMS4wNjU4IDQwLjkxOTMgMTEuODg1IDQxLjI0OTMgMTIuNzkxN0M0MS43MDI5IDE0LjAzNCA0MS45NDI1IDE1LjM0NDIgNDEuOTU3NyAxNi42NjY3QzQyLjA2MTggMTguNzUwIDQyLjA2MTggMTkuNSA0Mi4wNjE4IDI1QzQyLjA2MTggMzAuNSA0Mi4wNjE4IDMxLjI1IDQxLjk1NzcgMzMuMzMzNFoiIGZpbGw9IndoaXRlIi8+Cjwvc3ZnPg==" width="40" height="40"></a>
                                    </td>
                                    <td style="padding: 0 12px;">
                                        <a href="#"><img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNTAiIGhlaWdodD0iNTAiIHZpZXdCb3g9IjAgMCA1MCA1MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik00My45NTAxIDQzLjMzMzRMMjkuMjkxOCAyMS45NjM0TDI5LjMxNjggMjEuOTgzNEw0Mi41MzM1IDYuNjY2NjlIMzguMTE2OEwyNy4zNTAxIDE5LjEzMzRMMTguODAwMSA2LjY2NjY5SDcuMjE2OEwyMC45MDE4IDI2LjYxODRMMjAuOTAwMSAyNi42MTY3TDYuNDY2OCA0My4zMzM0SDEwLjg4MzVMMjIuODUzNSAyOS40NjM0TDMyLjM2NjggNDMuMzMzNEg0My45NTAxWk0xNy4wNTAxIDEwTDM3LjYxNjggNDBIMzQuMTE2OEwxMy41MzM1IDEwSDE3LjA1MDFaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4=" width="40" height="40"></a>
                                    </td>
                                    <td style="padding: 0 12px;">
                                        <a href="#"><img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNTAiIGhlaWdodD0iNTAiIHZpZXdCb3g9IjAgMCA1MCA1MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik00Ny45MTYyIDIwLjIyOTJDNDguMDE5MiAxNy4yNDcyIDQ3LjM2NzEgMTQuMjg3NiA0Ni4wMjAzIDExLjYyNTBDNDUuMTA2NiAxMC41MzI1IDQzLjgzODUgOS43OTUyNiA0Mi40MzcgOS41NDE2OUMzNi42NDAyIDkuMDE1NzEgMzAuODE5MyA4LjgwMDEyIDI0Ljk5OTUgOC44OTU4NkMxOS4yMDA4IDguNzk1NzcgMTMuNDAxIDkuMDA0NDAgNy42MjQ0OSA5LjUyMDg2QzYuNDgyNDQgOS43Mjg2IDUuNDI1NTcgMTAuMjY0MyA0LjU4MjgyIDExLjA2MjVDMi43MDc4MiAxMi43OTE3IDIuNDk5NDkgMTUuNzUwIDIuMjkxMTUgMTguMjUwQzEuOTg4ODkgMjIuNzQ1IDEuOTg4ODkgMjcuMjU1MSAyLjI5MTE1IDMxLjc1MEMyLjM1MTQyIDMzLjE1NzEgMi41NjA5MyAzNC41NTM4IDIuOTE2MTUgMzUuOTE2N0MzLjE2NzM1IDM2Ljk2ODkgMy42NzU1OCAzNy45NDI0IDQuMzk1MzIgMzguNzUwQzUuMjQzNzkgMzkuNTkwNiA2LjMyNTI5IDQwLjE1NjcgNy40OTk0OSA0MC4zNzVDMTEuOTkxIDQwLjkyOTQgMTYuNTE2NiA0MS4xNTkyIDIxLjA0MTIgNDEuMDYyNUMyOC4zMzI4IDQxLjE2NjcgMzQuNzI4NyA0MS4wNjI1IDQyLjI5MTIgNDAuNDc5MkM0My40OTQyIDQwLjI3NDMgNDQuNjA2MSAzOS43MDc0IDQ1LjQ3ODcgMzguODU0MkM0Ni4wNjE5IDM4LjI3MDcgNDYuNDk3NiAzNy41NTY1IDQ2Ljc0OTUgMzYuNzcwOUM0Ny40OTQ1IDM0LjQ4NDYgNDcuODYwNSAzMi4wOTE5IDQ3LjgzMjggMjkuNjg3NUM0Ny45MTYyIDI4LjUyMDkgNDcuOTE2MiAyMS40NzkyIDQ3LjkxNjIgMjAuMjI5MlpNMjAuMjkxMiAzMC45Mzc1VjE4LjA0MTdMMzIuNjI0NSAyNC41MjA5QzI5LjE2NjIgMjYuNDM3NSAyNC42MDM3IDI4LjYwNDIgMjAuMjkxMiAzMC45Mzc1WiIgZmlsbPSJ3aGl0ZSIvPgo8L3N2Zz4=" width="40" height="40"></a>
                                    </td>
                                </tr>
                            </table>

                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</div>
	`, otp)
}
