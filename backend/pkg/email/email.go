package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	apiKey    string
	fromEmail string
	fromName  string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey:    os.Getenv("SENDGRID_API_KEY"),
		fromEmail: os.Getenv("FROM_EMAIL"),
		fromName:  "10000 Hour Tracker",
	}
}

func (s *EmailService) SendEmail(toEmail, subject, body string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, body, body)

	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %s", response.Body)
	}

	return nil
}

func (s *EmailService) SendOTP(toEmail, code string) error {
	subject := "Your Verification Code - 10000 Hour Tracker"
	body := fmt.Sprintf(`
Hello,

Your verification code is: %s

This code will expire in 10 minutes.

If you didn't request this code, please ignore this email.

Best regards,
10000 Hour Tracker Team
	`, code)

	return s.SendEmail(toEmail, subject, body)
}
