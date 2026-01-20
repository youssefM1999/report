package mailer

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewEmailData(t *testing.T) {
	subject := "Weekly Report"
	markdown := `## social

### 8ae1b21 - Implement User Validation
- Added migration to create user_invitations table
- Implemented email activation mechanism
`
	period := 7 * 24 * time.Hour

	data := NewEmailData(subject, markdown, period)

	if data.Subject != subject {
		t.Errorf("Subject = %q, want %q", data.Subject, subject)
	}
	if data.Period != "Past Week" {
		t.Errorf("Period = %q, want %q", data.Period, "Past Week")
	}
	if len(data.HTMLContent) == 0 {
		t.Error("HTMLContent should not be empty")
	}
	if len(data.GeneratedAt) == 0 {
		t.Error("GeneratedAt should not be empty")
	}
}

func TestRenderEmailTemplate(t *testing.T) {
	markdownReport := `## social

### 8ae1b21 - Implement User Validation
- Added migration to create user_invitations table with token and expiry columns
- Implemented email activation mechanism with backend token validation
- Updated user registration flow to include invitation process

### 1df40a6 - Improve email functionality
- Enhanced mailer package with more robust error handling
- Updated email sending process to log detailed information
- Implemented retry mechanism for email sending

## report

No commits in this period.
`

	data := NewEmailData("Developer Activity Report", markdownReport, 7*24*time.Hour)

	html, err := renderEmailTemplate(data)
	if err != nil {
		t.Fatalf("renderEmailTemplate() failed: %v", err)
	}

	// Write to file for preview
	if err := os.WriteFile("/tmp/report_preview.html", []byte(html), 0644); err != nil {
		t.Logf("Could not write preview file: %v", err)
	} else {
		t.Log("Preview saved to /tmp/report_preview.html - open in browser to view")
	}

	// Basic checks
	if len(html) == 0 {
		t.Error("HTML output is empty")
	}
	if !strings.Contains(html, "Developer Activity Report") {
		t.Error("HTML should contain title")
	}
	if !strings.Contains(html, "8ae1b21") {
		t.Error("HTML should contain commit hash")
	}
	if !strings.Contains(html, "Past Week") {
		t.Error("HTML should contain period")
	}
	if !strings.Contains(html, "social") {
		t.Error("HTML should contain repo name")
	}
}

func TestMarkdownToHTML(t *testing.T) {
	md := `## Repository Name

### abc1234 - Fix bug
- Fixed the thing
- Updated tests
`
	html := markdownToHTML(md)

	if !strings.Contains(html, "<h2") {
		t.Error("Should contain h2 tag")
	}
	if !strings.Contains(html, "<h3") {
		t.Error("Should contain h3 tag")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("Should contain li tags")
	}

	t.Log("Generated HTML:")
	t.Log(html)
}

func TestFormatPeriod(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{7 * 24 * time.Hour, "Past Week"},
		{14 * 24 * time.Hour, "Past 2 Weeks"},
		{30 * 24 * time.Hour, "Past Month"},
		{5 * 24 * time.Hour, "Past 5 Days"},
	}

	for _, tt := range tests {
		result := formatPeriod(tt.duration)
		if result != tt.expected {
			t.Errorf("formatPeriod(%v) = %q, want %q", tt.duration, result, tt.expected)
		}
	}
}
