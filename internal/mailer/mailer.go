package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	FromName   = "Report My Work 4 Me"
	FromEmail  = "reportmywork4me@gmail.com"
	MaxRetries = 3
)

type Client interface {
	Send(email, username, subject, markdownContent string, period time.Duration, isSandbox bool) error
}

type EmailData struct {
	Subject     string
	HTMLContent template.HTML
	Period      string
	GeneratedAt string
}

func NewEmailData(subject, markdownContent string, period time.Duration) EmailData {
	return EmailData{
		Subject:     subject,
		HTMLContent: template.HTML(markdownToHTML(markdownContent)),
		Period:      formatPeriod(period),
		GeneratedAt: time.Now().Format("January 2, 2006 at 3:04 PM"),
	}
}

func markdownToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func formatPeriod(d time.Duration) string {
	days := int(d.Hours() / 24)
	switch days {
	case 7:
		return "Past Week"
	case 14:
		return "Past 2 Weeks"
	case 30:
		return "Past Month"
	default:
		return fmt.Sprintf("Past %d Days", days)
	}
}

const emailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        h2 { color: #1a1a1a; font-size: 18px; margin-top: 25px; margin-bottom: 12px; padding-bottom: 6px; border-bottom: 1px solid #e1e5e9; }
        h3 { color: #2c5282; font-size: 15px; margin-top: 18px; margin-bottom: 8px; }
        ul { margin: 8px 0; padding-left: 24px; }
        li { margin-bottom: 5px; color: #444; }
        code { background: #f5f5f5; padding: 1px 5px; border-radius: 3px; font-size: 13px; }
    </style>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px;">

    <div style="border-bottom: 2px solid #2c5282; padding-bottom: 15px; margin-bottom: 25px;">
        <h1 style="margin: 0; font-size: 24px; font-weight: 600; color: #1a1a1a;">Developer Activity Report</h1>
        <p style="margin: 5px 0 0 0; color: #666; font-size: 14px;">{{.Period}} • Generated {{.GeneratedAt}}</p>
    </div>

    <div style="font-size: 15px;">
        {{.HTMLContent}}
    </div>

    <div style="margin-top: 30px; padding-top: 15px; border-top: 1px solid #e1e5e9; font-size: 12px; color: #888;">
        <p style="margin: 0;">This report was automatically generated. Feel free to edit before forwarding.</p>
    </div>

</body>
</html>`

func renderEmailTemplate(data EmailData) (string, error) {
	// Add some inline styles for the markdown-generated HTML
	styledTemplate := emailTemplate

	tmpl, err := template.New("email").Parse(styledTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
