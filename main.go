package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// SMTPSettings struct holds the configuration for the SMTP server
type SMTPSettings struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// Mailer struct to manage email sending
type Mailer struct {
	settings SMTPSettings // Embed SMTPSettings within Mailer
}

// EmailPayload struct to manage the email and metadata itself
type EmailPayload struct {
	to             string
	subject        string
	body           string
	attachmentPath string
	yearsAway      int
}

// NewMailer constructor function for creating a new Mailer instance
func NewMailer(settings SMTPSettings) *Mailer {
	return &Mailer{
		settings: settings,
	}
}

// SendEmailWithAttachmentOnDate checks if the current date is given amount of years from burn, if matches payload, returns true
func (m *Mailer) CalculateDateDiff(payload *EmailPayload) (bool, error) {
	burnDate := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC) // date of the burn
	currentDate := time.Now()
	yearsDiff := currentDate.Year() - burnDate.Year()
	if yearsDiff != payload.yearsAway {
		// If not the specified years away, return an error
		return false, fmt.Errorf("current date is not %d years away from the set date", payload.yearsAway)
	} else {
		return true, nil
	}
}

func parseEmailPayload(dirPath, filename string) (EmailPayload, error) {
	parts := strings.Split(filename, "_")
	if len(parts) != 2 {
		return EmailPayload{}, fmt.Errorf("invalid filename format: %s", filename)
	}
	email := parts[0]

	yearsAndFiletype := strings.Split(parts[1], ".")
	if len(yearsAndFiletype) != 2 {
		return EmailPayload{}, fmt.Errorf("invalid filename format: %s", filename)
	}

	years, err := strconv.Atoi(yearsAndFiletype[0])
	if err != nil {
		return EmailPayload{}, fmt.Errorf("unable to parse years from filename: %s", filename)
	}

	emailPayload := EmailPayload{
		to:             email,
		subject:        "Hello from Eternal Draft",
		body:           "TEST BODY",
		attachmentPath: dirPath + filename,
		yearsAway:      years,
	}

	return emailPayload, nil
}

func listPostcards(dirPath string) ([]EmailPayload, error) {
	var payloads []EmailPayload

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
		// Read the attachment file
	}

	for _, file := range files {
		if !file.IsDir() {
			// parse the email payloads here
			filename := file.Name()
			emailPayload, err := parseEmailPayload(dirPath, filename)
			if err != nil {
				log.Printf("[!] Error parsing E-Mail payload for file %s: %v", filename, err)
				continue // skipping this file and continuing
			}
			payloads = append(payloads, emailPayload)
		}

	}
	return payloads, nil
}

// sendEmailWithAttachment handles the email sending with an image attachment
func (m *Mailer) sendEmailWithAttachment(payload *EmailPayload) error {

	fileBytes, err := os.ReadFile(payload.attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %w", err)
	}

	// Create a buffer to build the email message
	var email bytes.Buffer
	writer := multipart.NewWriter(&email)

	// Write the email headers (from, to, subject)
	fmt.Fprintf(&email, "From: %s\r\n", m.settings.From)
	fmt.Fprintf(&email, "To: %s\r\n", payload.to)
	fmt.Fprintf(&email, "Subject: %s\r\n", payload.subject)
	fmt.Fprintf(&email, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&email, "Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary())
	email.Write([]byte("\r\n")) // End of headers section

	// Add the email body as a part
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to add body part: %w", err)
	}
	part.Write([]byte(payload.body + "\r\n"))

	// Add the attachment as a part
	part, err = writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition":       []string{fmt.Sprintf("attachment; filename=\"%s\"", payload.attachmentPath)},
		"Content-Type":              []string{"application/octet-stream"},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return fmt.Errorf("failed to add attachment part: %w", err)
	}
	b64 := base64.NewEncoder(base64.StdEncoding, part)
	if _, err := b64.Write(fileBytes); err != nil {
		return fmt.Errorf("failed to encode attachment: %w", err)
	}
	b64.Close()

	// Close the multipart writer to finish writing the email body
	writer.Close()

	// Send the email
	addr := fmt.Sprintf("%s:%d", m.settings.Host, m.settings.Port)
	auth := smtp.PlainAuth("", m.settings.Username, m.settings.Password, m.settings.Host)
	if err := smtp.SendMail(addr, auth, m.settings.From, []string{payload.to}, email.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Printf("[*] Successfully sent mail to %s", payload.to)

	return nil
}

func main() {

	// TODO Google SMTP smtp-relay

	// Setup SMTP
	fmt.Println("[*] Setting up...")
	settings := SMTPSettings{
		Host:     os.Getenv("SMPT_HOST"),
		Port:     587,
		From:     os.Getenv("MAIL_FROM"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
	// Create a new mailer instance
	mailer := NewMailer(settings)
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("[!] %+v", err)
		panic("[!] Failed to get current working directory")
	}

	dirPath := filepath.Join(cwd, "postcards")
	payloads, err := listPostcards(dirPath)
	if err != nil {
		fmt.Println(err)
		panic("A problem!")
	}

	for _, payload := range payloads {
		dateDiff, err := mailer.CalculateDateDiff(&payload)
		if err != nil {
			log.Printf("[!] %+v for %+v", err, payload)
		}

		if dateDiff {
			err := mailer.sendEmailWithAttachment(&payload)
			if err != nil {
				fmt.Printf("[ERROR] Failed to send email: %v\n", err)
			}
		}
	}
}
