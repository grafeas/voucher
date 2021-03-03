package subscriber

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/grafeas/voucher/cmd/config"
	"github.com/grafeas/voucher/metrics"
	"github.com/sirupsen/logrus"
)

// Subscriber contains the information required to pull messages from a pub/sub topic.
type Subscriber struct {
	cfg     *Config
	secrets *config.Secrets
	metrics metrics.Client
	log     *logrus.Logger
}

// NewSubscriber creates a new subscription topic puller for a subscription.
func NewSubscriber(config *Config, secrets *config.Secrets, metrics metrics.Client, log *logrus.Logger) *Subscriber {
	return &Subscriber{
		cfg:     config,
		secrets: secrets,
		metrics: metrics,
		log:     log,
	}
}

// Subscribe pulls messages for a subscription and passes them along to get checked.
func (s *Subscriber) Subscribe(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, s.cfg.Project)
	if err != nil {
		return fmt.Errorf("failed to create new pubsub client: %s", err)
	}
	defer client.Close()

	sub := client.Subscription(s.cfg.Subscription)
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		processStart := time.Now()
		defer func(startTime time.Time) {
			s.metrics.PubSubTotalLatency(time.Since(startTime))
		}(processStart)

		s.metrics.PubSubMessageReceived()

		pl, err := parsePayload(msg.Data)
		if err != nil {
			if err != errNotInsertAction {
				s.log.WithField("reason", err).WithField("payload", string(msg.Data)).Error("couldn't parse pub/sub payload")
			}

			msg.Ack()
			return
		}

		l := s.log.WithField("payload", pl)

		cir, err := pl.asCanonicalImage()
		if err != nil {
			msg.Ack()
			l.WithField("reason", err).Error("couldn't make canonical image")
			return
		}

		l.WithField("status", "pending").Info("the vouch started")
		vouchStatus, shouldRetry := s.check(cir)

		// Nack if we want to retry, otherwise the catch-all Ack below will
		// ensure this message doesn't get retried
		if vouchStatus {
			l.WithField("status", "success").Info("the vouch succeeded")
		} else {
			l.WithField("status", "failure").Info("the vouch failed")

			if shouldRetry {
				msg.Nack()
			}
		}

		msg.Ack()
	})

	if err != nil {
		return fmt.Errorf("sub.Receive: %s", err)
	}

	return nil
}
