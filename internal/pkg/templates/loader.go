package templates

import (
	"fmt"
	"os"
	"time"
)

const (
	emailCodeTemplate       = "./templates/email_code.html"
	invoiceCompleteTemplate = "./templates/invoice_complete.html"
	welcomeTemplate         = "./templates/welcome.html"
)

func load(templatePath string) (string, error) {
	htmlBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Failed to read HTML file:", err)
		return "", fmt.Errorf("os.ReadFile: %w", err)
	}

	return string(htmlBytes), nil
}

func BuildEmailCode(emailCode string) (string, error) {
	content, err := load(emailCodeTemplate)
	if err != nil {
		return "", fmt.Errorf("load: %w", err)
	}

	// We need you time to prevent message being trimmed
	return fmt.Sprintf(content, emailCode, time.Now()), nil
}

func BuildInvoiceComplete(invoiceID string, usdAmount float64) (string, error) {
	content, err := load(invoiceCompleteTemplate)
	if err != nil {
		return "", fmt.Errorf("load: %w", err)
	}

	// We need you time to prevent message being trimmed
	return fmt.Sprintf(content, invoiceID, usdAmount, time.Now()), nil
}

func BuildWelcome(username string) (string, error) {
	content, err := load(welcomeTemplate)
	if err != nil {
		return "", fmt.Errorf("load: %w", err)
	}

	// We need you time to prevent message being trimmed
	return fmt.Sprintf(content, username, time.Now()), nil
}
