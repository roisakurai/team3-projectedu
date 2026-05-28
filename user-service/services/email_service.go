package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendWelcomeEmail(toEmail string, toName string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	if apiKey == "" || fromEmail == "" {
		return errors.New("sendgrid configuration is missing")
	}

	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	subject := "Welcome to EduAccess"

	plainTextContent := fmt.Sprintf(
		"Hi %s, welcome to EduAccess. Your account has been successfully registered.",
		toName,
	)

	htmlContent := fmt.Sprintf(`
<h2>Welcome to EduAccess %s</h2>

<p>
Thank you for registering your account.
</p>

<p>
You can now access classes, learning materials,
assignments, and track your progress.
</p>

<p>
If you did not create this account,
please ignore this email.
</p>
	`, toName)

	message := mail.NewSingleEmail(
		from,
		subject,
		to,
		plainTextContent,
		htmlContent,
	)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

func (s *EmailService) SendVerificationEmail(
	toEmail string,
	toName string,
	token string,
) error {

	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	fromName := os.Getenv("SENDGRID_FROM_NAME")

	verificationLink := fmt.Sprintf(
		"http://localhost:8080/verify-email/%s",
		token,
	)

	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	subject := "Verify Your EduAccess Account"

	plainTextContent := fmt.Sprintf(
		"Click this link to verify your account: %s",
		verificationLink,
	)

	htmlContent := fmt.Sprintf(`
		<h2>Verify Your Email</h2>

		<p>Hello %s,</p>

		<p>
			Thank you for registering at EduAccess.
		</p>

		<p>
			Please click the button below to verify your account:
		</p>

		<a href="%s"
		   style="
				background:#2563eb;
				color:white;
				padding:12px 20px;
				text-decoration:none;
				border-radius:8px;
				display:inline-block;
		   ">
			Verify Email
		</a>

		<p>
			This link will expire in 24 hours.
		</p>
	`, toName, verificationLink)

	message := mail.NewSingleEmail(
		from,
		subject,
		to,
		plainTextContent,
		htmlContent,
	)

	client := sendgrid.NewSendClient(apiKey)

	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf(
			"failed to send email: %d %s",
			response.StatusCode,
			response.Body,
		)
	}

	return nil
}
