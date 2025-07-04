package mail

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"
)

func SendNotificationMail(subject, message string) error {
	// Get API key from environment variable
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Printf("RESEND_API_KEY environment variable not set")
		return nil // Return nil to not break the application if email is not configured
	}

	// Get recipient emails from environment variable (comma-separated)
	toEmails := os.Getenv("NOTIFICATION_EMAILS")
	if toEmails == "" {
		log.Printf("NOTIFICATION_EMAILS environment variable not set")
		return nil // Return nil to not break the application if email is not configured
	}

	// Split comma-separated emails and trim whitespace
	emailList := make([]string, 0)
	for _, email := range strings.Split(toEmails, ",") {
		trimmedEmail := strings.TrimSpace(email)
		if trimmedEmail != "" {
			emailList = append(emailList, trimmedEmail)
		}
	}

	if len(emailList) == 0 {
		log.Printf("No valid email addresses found in NOTIFICATION_EMAILS")
		return nil
	}

	// Check network connectivity to Resend API
	_, err := net.DialTimeout("tcp", "api.resend.com:443", 5*time.Second)
	if err != nil {
		log.Printf("Network connectivity issue - cannot reach Resend API: %v", err)
		log.Printf("Email notification skipped for: %s", subject)
		return nil // Return nil to not break the application
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      emailList,
		Subject: subject,
		Html:    "<p>" + message + "</p>",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		log.Printf("Email notification skipped for: %s", subject)
		return nil // Return nil instead of error to not break the application
	}

	log.Printf("Email sent successfully: %v", sent)
	return nil
}
