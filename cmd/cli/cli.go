package cli

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	emailFlag = &cli.StringFlag{
		Name:  "email",
		Usage: "The email address to send the report to",
	}
	rangeFlag = &cli.DurationFlag{
		Name:  "range",
		Usage: "The range of time to generate the report for",
		Value: 7 * 24 * time.Hour,
	}
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "report",
		Usage: "Generate a report of your work",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate a report of your work",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "email",
						Usage: "The email address to send the report to",
					},
				},
			},
		},
	}
}

func generateReport() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate a report of your work",
		Flags: []cli.Flag{
			emailFlag,
			rangeFlag,
		},
	}
}
