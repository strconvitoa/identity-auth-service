package services

import (
	"fmt"
	"net/smtp"

	"github.com/strconvitoa/martian-service/internal/core/ports"
)

type emailService struct {
}

func NewEmailService() ports.EmailService {
	return &emailService{}
}

func (s *emailService) SendEmail(toEmail string, subject string, body string) error {
	// 1. SMTP Configuration
	email := "admin@martian.esq"
	from := "crew@martian.esq"
	password := "slxdqvsbbzhbzsak"
	host := "smtp.gmail.com"
	port := "587"
	to := []string{toEmail}

	// 1. Build the headers correctly
	// The "Subject:" prefix is mandatory for the subject to show up.
	// The "From:" header tells the recipient's email client to show the alias.
	headerFrom := "From: " + from + "\r\n"
	headerTo := "To: " + toEmail + "\r\n"
	headerSubject := "Subject: " + subject + "\r\n"
	mime := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" // This empty line is CRITICAL to separate headers from the body

	// 2. Combine headers and body
	msg := []byte(headerFrom + headerTo + headerSubject + mime + body)

	// 3. Authenticate as the PRIMARY account (email/password)
	auth := smtp.PlainAuth("", email, password, host)

	// 4. Send as the ALIAS (from)
	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *emailService) GetVerificationEmailBody(otpCode string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500&family=JetBrains+Mono:wght@500&display=swap');
        </style>
    </head>
    <body style="margin:0; padding:0; background-color:#f9f9f9; font-family:'Inter', sans-serif; color:#1a1c1c;">
        <div style="max-width:600px; margin:40px auto; padding:20px;">
            <div style="background-color:#ffffff; border:1px solid #D1D3D5; border-radius:4px; overflow:hidden;">
                <div style="padding:40px; border-bottom:1px solid #f0f0f0; display:flex; justify-content:between;">
                    <span style="font-weight:bold; font-size:18px; color:#2C3539;">
                    <img src="../../../img/logo.png" alt="logo" style="vertical-align: middle; height: 1em; margin-right: 8px;">
                    </span>
                    <span style="font-size:11px; font-weight:bold; letter-spacing:0.1em; color:#2C3539; text-transform:uppercase; margin-left:auto;">ACCOUNT INFORMATION</span>
                </div>
                <div style="padding:48px 40px;">
                    <h1 style="font-size:32px; font-weight:500; margin-bottom:24px; color:#1a1c1c;">Welcome, Martian </h1>
                    <p style="font-size:16px; line-height:1.5; color:#444749; margin-bottom:40px;">Please use the following verification code to log in</p>
                    <div style="background-color:#f3f3f4; border-radius:4px; padding:40px; text-align:center;">
                        <span style="font-size:48px; font-weight:600; letter-spacing:-0.02em; color:#2C3539;">%s</span>
                    </div>
                </div>
            </div>
            <div style="margin-top:40px; text-align:center; color:#99A2A8; font-size:10px; font-family:'JetBrains Mono', monospace;">
                <p>25 Kent Ave, Brooklyn, NY 11249</p>
                <p>spread love ❤️ is the brooklyn way</p>
            </div>
        </div>
    </body>
    </html>`, otpCode)
}
func (s *emailService) GetWelcomeEmailBody(otpCode string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500&family=JetBrains+Mono:wght@500&display=swap');
        </style>
    </head>
    <body style="margin:0; padding:0; background-color:#f9f9f9; font-family:'Inter', sans-serif; color:#1a1c1c;">
        <div style="max-width:600px; margin:40px auto; padding:20px;">
            <div style="background-color:#ffffff; border:1px solid #D1D3D5; border-radius:4px; overflow:hidden;">
                <div style="padding:40px; border-bottom:1px solid #f0f0f0; display:flex; justify-content:between;">
                    <span style="font-weight:bold; font-size:18px; color:#2C3539;">
                    <img src="https://127.0.0.1:8443/static/img/logo.png" alt="martian" style="vertical-align: middle; height: 40px; width: auto; margin-right: 8px;">
                    </span>
                    <span style="font-size:11px; font-weight:bold; letter-spacing:0.1em; color:#2C3539; text-transform:uppercase; margin-left:auto;">Confirm your email address</span>
                </div>
                <div style="padding:48px 40px;">
                    <h1 style="font-size:32px; font-weight:500; margin-bottom:24px; color:#1a1c1c;">Verification Code</h1>
                    <p style="font-size:16px; line-height:1.5; color:#444749; margin-bottom:40px;">Thank you! Please enter the code below in your open web browser window or click the button below to verify your email.</p>
                    <div style="background-color:#f3f3f4; border-radius:4px; padding:40px; text-align:center;">
                        <span style="font-size:48px; font-weight:600; letter-spacing:-0.02em; color:#2C3539;">%s</span>
                    </div>
                    <p style="font-size:16px; line-height:1.5; color:#444749; margin-bottom:40px;">You’re receiving this email because your account sign-up process with martian.esq is not yet complete. If you no longer wish to create an account, you can safely ignore this email.</p>
                </div>
            </div>
            <div style="margin-top:40px; text-align:center; color:#99A2A8; font-size:10px; font-family:'JetBrains Mono', monospace;">
                <p>Martian.esq, 25 Kent Ave, Brooklyn, NY 11249</p>
                <p>spread love ❤️ is the brooklyn way</p>
            </div>
        </div>
    </body>
    </html>`, otpCode)
}
