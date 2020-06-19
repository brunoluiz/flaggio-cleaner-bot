package worker

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/brunoluiz/flaggio-cleaner-bot/flaggio"
	"github.com/brunoluiz/flaggio-cleaner-bot/linear"
	"github.com/brunoluiz/flaggio-cleaner-bot/repo"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// FlaggioClient flags repo
type FlaggioClient interface {
	FindFlagsByMaxAge(ctx context.Context, maxAge time.Duration, opts ...flaggio.FindFlagsOpt) ([]flaggio.Flag, error)
}

// ProcessedFlagsRepository processed flags repo
type ProcessedFlagsRepository interface {
	IsProcessed(ctx context.Context, id string) (bool, error)
	Save(ctx context.Context, flag repo.ProcessedFlag) error
}

// LinearClient linear client
type LinearClient interface {
	CreateIssue(ctx context.Context, in linear.IssueCreateInput) (linear.IssuePayload, error)
}

// Worker defines a worker
type Worker struct {
	flags     FlaggioClient
	processed ProcessedFlagsRepository
}

// New return a new worker
func New(
	flags FlaggioClient,
	processed ProcessedFlagsRepository,
) *Worker {
	return &Worker{flags, processed}
}

// ProcessOutdatedTrigger Triggers for process outdated
type ProcessOutdatedTrigger func(ctx context.Context, f flaggio.Flag) error

// WithLinearIssueCreation Creates an issue at linear if a flag is outdated
func WithLinearIssueCreation(
	client LinearClient,
	teamID string,
	projectID string,
	templateStr string,
) (ProcessOutdatedTrigger, error) {
	t, err := template.New("linear").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, f flaggio.Flag) error {
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, f); err != nil {
			return err
		}

		_, err := client.CreateIssue(ctx, linear.IssueCreateInput{
			Title:     tpl.String(),
			TeamID:    teamID,
			ProjectID: projectID,
		})
		return err
	}, nil
}

// WithSlackNotification Notify slack if a flag is outdated
func WithSlackNotification(
	client *slack.Client,
	channelID string,
	templateStr string,
) (ProcessOutdatedTrigger, error) {
	t, err := template.New("slack").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, f flaggio.Flag) error {
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, f); err != nil {
			return err
		}

		_, _, err := client.PostMessageContext(
			ctx,
			channelID,
			slack.MsgOptionText(tpl.String(), true),
		)
		return err
	}, nil
}

// ProcessOutdated scan old flags and create tickets on linear for clean-up
func (w *Worker) ProcessOutdated(
	ctx context.Context,
	maxAge time.Duration,
	flagPrefix string,
	triggers ...ProcessOutdatedTrigger,
) error {
	if len(triggers) == 0 {
		return errors.New("No trigger configured")
	}

	flags, err := w.flags.FindFlagsByMaxAge(ctx, maxAge, flaggio.WithFindFlagsSearchOpt(flagPrefix))
	if err != nil {
		return err
	}

	var errs error
	for _, f := range flags {
		isProcessed, err := w.processed.IsProcessed(ctx, f.ID)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if isProcessed {
			continue
		}

		var triggerErrs error
		for _, trigger := range triggers {
			if err := trigger(ctx, f); err != nil {
				triggerErrs = multierror.Append(triggerErrs, err)
				continue
			}
		}

		if triggerErrs != nil {
			errs = multierror.Append(errs, triggerErrs)
			continue
		}

		if err = w.processed.Save(ctx, repo.ProcessedFlag{
			ID:   f.ID,
			Key:  f.Key,
			Name: f.Name,
		}); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"flag": f.Key,
		}).Infof("Successfully processed feature flag '%s'", f.Key)
	}

	return errs
}
