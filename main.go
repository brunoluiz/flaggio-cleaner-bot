package main

import (
	"context"
	"os"
	"time"

	"github.com/brunoluiz/flaggio-cleaner-bot/flaggio"
	"github.com/brunoluiz/flaggio-cleaner-bot/linear"
	"github.com/brunoluiz/flaggio-cleaner-bot/repo"
	"github.com/brunoluiz/flaggio-cleaner-bot/worker"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "flaggio-cleaner-bot",
		Usage: "Process and notify old Flaggio feature flags 🤖",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "max-age",
				EnvVars: []string{"MAX_AGE"},
				Value:   time.Hour * 24 * 5,
				Usage:   "Max flag age before processing is triggered",
			},
			&cli.StringFlag{
				Name:     "storage-path",
				EnvVars:  []string{"STORAGE_PATH"},
				Usage:    "Used to store already processed flags",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "flaggio-url",
				EnvVars:  []string{"FLAGGIO_URL"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "flag-prefix",
				EnvVars: []string{"FLAG_PREFIX"},
				Value:   "",
				Usage:   "Searches for a certain flag prefix -- useful for shared flaggio instances",
			},
			&cli.StringFlag{
				Name:    "linear-token",
				EnvVars: []string{"LINEAR_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "linear-team",
				EnvVars: []string{"LINEAR_TEAM"},
			},
			&cli.StringFlag{
				Name:    "linear-project",
				EnvVars: []string{"LINEAR_PROJECT"},
			},
			&cli.StringFlag{
				Name:    "linear-template",
				EnvVars: []string{"LINEAR_TEMPLATE"},
				Value:   "Clean-up code for feature flag '{{ .Key }}'",
			},
			&cli.StringFlag{
				Name:    "slack-token",
				EnvVars: []string{"SLACK_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "slack-channel",
				EnvVars: []string{"SLACK_CHANNEL"},
			},
			&cli.StringFlag{
				Name:    "slack-template",
				EnvVars: []string{"SLACK_TEMPLATE"},
				Value:   "Time to clean-up code for feature flag `{{ .Key }}` 🔥",
			},
		},
		Action: func(c *cli.Context) error {
			return run(
				c.Context,
				c.Duration("max-age"),
				c.String("storage-path"),
				c.String("flaggio-url"),
				c.String("flag-prefix"),
				c.String("linear-token"),
				c.String("linear-team"),
				c.String("linear-project"),
				c.String("linear-template"),
				c.String("slack-token"),
				c.String("slack-channel"),
				c.String("slack-template"),
			)
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(
	ctx context.Context,
	maxAge time.Duration,
	storagePath string,
	flaggioURL string,
	flagPrefix string,
	linearToken string,
	linearTeam string,
	linearProject string,
	linearTemplate string,
	slackToken string,
	slackChannel string,
	slackTemplate string,
) error {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	processedFlag, err := repo.NewProcessedFlagDisk(storagePath)
	if err != nil {
		return err
	}
	defer processedFlag.Close()

	w := worker.New(
		flaggio.New(flaggioURL),
		processedFlag,
	)

	triggers := []worker.ProcessOutdatedTrigger{}
	if linearToken != "" && linearTeam != "" {
		t, err := worker.WithLinearIssueCreation(
			linear.New(linearToken),
			linearTeam,
			linearProject,
			linearTemplate,
		)
		if err != nil {
			return err
		}

		triggers = append(triggers, t)
	}

	if slackToken != "" && slackChannel != "" {
		t, err := worker.WithSlackNotification(
			slack.New(slackToken),
			slackChannel,
			slackTemplate,
		)
		if err != nil {
			return err
		}

		triggers = append(triggers, t)
	}

	return w.ProcessOutdated(ctx, maxAge, flagPrefix, triggers...)
}
