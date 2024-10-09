package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Llane00/ramen-backend/initializers"
	"github.com/Llane00/ramen-backend/models"
	"github.com/k3a/html2text"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

type EmailRequest struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To []struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

const MailtrapApiUrl = "https://send.api.mailtrap.io/api/send"
const MailtrapTestApiUrl = "https://sandbox.api.mailtrap.io/api/send/3171336"

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.User, data *EmailData, emailTemp string) error {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data
	from := config.EmailFrom
	to := user.Email
	apiToken := config.SMTPApiToken // use API Token

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, emailTemp, &data)

	// Create email request
	emailReq := EmailRequest{
		From: struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			Email: from,
			Name:  "Sender Name",
		},
		To: []struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			{
				Email: to,
				Name:  user.Name,
			},
		},
		Subject: data.Subject,
		Text:    html2text.HTML2Text(body.String()),
		HTML:    body.String(),
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", MailtrapTestApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set request headers
	req.Header.Add("Authorization", "Bearer"+" "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email. Status: %s, Response: %s", resp.Status, string(respBody))
	}

	fmt.Println("Email sent successfully")
	return nil
}
