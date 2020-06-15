package main

import (
	"context"
	"os"
	"time"

	"github.com/brunoluiz/flaggio-cleaner-bot/internal/linear"
	"github.com/brunoluiz/flaggio-cleaner-bot/internal/repo"
	"github.com/brunoluiz/flaggio-cleaner-bot/internal/worker"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			},
			&cli.StringFlag{
				Name:     "db-dsn",
				EnvVars:  []string{"DB_DSN"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "flag-prefix",
				EnvVars: []string{"FLAG_PREFIX"},
				Value:   "",
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
			return run(c.Context,
				c.Duration("max-age"),
				c.String("db-dsn"),
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
	dbDSN string,
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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbDSN))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	w := worker.New(
		repo.NewFlagMongo(client.Database("flaggio")),
		repo.NewProcessedFlagMongo(client.Database("flaggio")),
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
