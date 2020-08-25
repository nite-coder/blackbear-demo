package service

import (
	"context"
	"sync"
	"time"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
)

type EventService struct {
	mu     sync.RWMutex
	config config.Configuration
	events []*event.Event
}

func NewEventService(cfg config.Configuration) event.EventServicer {
	return &EventService{
		config: cfg,
		events: []*event.Event{
			{
				ID:              1,
				Title:           "Golang Summit",
				Description:     "my desc",
				PublishedStatus: event.Draft,
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

// Events returns all events
func (srv *EventService) Events(ctx context.Context) ([]*event.Event, error) {
	return srv.events, nil
}

// UpdatePublishStatus update
func (srv *EventService) UpdatePublishStatus(ctx context.Context, request event.UpdateEventStatusRequest) error {
	logger := log.FromContext(ctx)
	logger.Debug("begin UpdatePublishStatus fn")

	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, evt := range srv.events {
		if evt.ID == request.EventID {
			evt.PublishedStatus = request.PublishedStatus
			evt.CreatedAt = time.Now().UTC()
		}
	}

	return nil
}
