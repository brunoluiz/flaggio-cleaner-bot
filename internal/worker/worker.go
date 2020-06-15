package worker

import (
	"bytes"
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/brunoluiz/flaggio-cleaner-bot/internal/linear"
	"github.com/brunoluiz/flaggio-cleaner-bot/internal/repo"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FlagRepository flags repo
type FlagRepository interface {
	FindFlagsByMaxAge(ctx context.Context, maxAge time.Duration) ([]repo.Flag, error)
}

// ProcessedFlagsRepository outdated flags repo
type ProcessedFlagsRepository interface {
	IsProcessed(ctx context.Context, id primitive.ObjectID) (bool, error)
	CreateProcessedFlag(ctx context.Context, flag repo.ProcessedFlag) error
}

// LinearClient linear client
type LinearClient interface {
	CreateIssue(ctx context.Context, in linear.IssueCreateInput) (linear.IssuePayload, error)
}

// Worker defines a worker
type Worker struct {
	flags     FlagRepository
	processed ProcessedFlagsRepository
}

// New return a new worker
func New(
	flags FlagRepository,
	outdated ProcessedFlagsRepository,
) *Worker {
	return &Worker{flags, outdated}
}

// ProcessOutdatedTrigger Triggers for process outdated
type ProcessOutdatedTrigger func(ctx context.Context, f repo.Flag) error

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

	return func(ctx context.Context, f repo.Flag) error {
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

	return func(ctx context.Context, f repo.Flag) error {
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

	flags, err := w.flags.FindFlagsByMaxAge(ctx, maxAge)
	if err != nil {
		return err
	}

	var errs error
	for _, f := range flags {
		hasPrefix := strings.HasPrefix(f.Key, flagPrefix)
		if !hasPrefix {
			continue
		}

		isNotified, err := w.processed.IsProcessed(ctx, f.ID)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if isNotified {
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

		if err = w.processed.CreateProcessedFlag(ctx, repo.ProcessedFlag{
			ID:  f.ID,
			Key: f.Key,
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
