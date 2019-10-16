package eventbridge_logrus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/sirupsen/logrus"
)

type eventbrigeHook struct {
	svc      *eventbridge.Client
	source   string
	eventBus string
}

func NewEventbridgeHook(region, source, eventBus string) (*eventbrigeHook, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Region = region
	svc := eventbridge.New(cfg)
	return &eventbrigeHook{
		svc:      svc,
		source:   source,
		eventBus: eventBus,
	}, nil
}

func (h *eventbrigeHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.JSONFormatter{}
	ent, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	req := h.svc.PutEventsRequest(&eventbridge.PutEventsInput{
		Entries: []eventbridge.PutEventsRequestEntry{
			{
				Detail:       aws.String(string(ent)),
				DetailType:   aws.String("Log messages from logrus"),
				EventBusName: aws.String(h.eventBus),
				Source:       aws.String(h.source),
				Time:         aws.Time(entry.Time),
			}},
	})
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	resp, err := req.Send(ctx)
	if resp.FailedEntryCount != nil && *resp.FailedEntryCount > 0 {
		ent := resp.Entries[0]
		return errors.New(fmt.Sprintf("%s - %s", *ent.ErrorCode, *ent.ErrorMessage))
	}
	return err
}

func (h *eventbrigeHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
