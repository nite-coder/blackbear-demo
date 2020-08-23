package service

import (
	"context"
	"testing"

	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	ctx := context.Background()
	cfg := config.Configuration{}
	eventService := NewEventService(cfg)

	events, err := eventService.Events(ctx)
	require.NoError(t, err, "list events should be no errors")

	assert.Equal(t, 1, len(events))

	evt := events[0]
	assert.Equal(t, int64(1), evt.ID)
	assert.Equal(t, event.Draft, evt.PublishedStatus)
}

func TestUpdatePublishStatus(t *testing.T) {
	ctx := context.Background()
	cfg := config.Configuration{}
	eventService := NewEventService(cfg)

	request := event.UpdateEventStatusRequest{
		EventID:         1,
		TransID:         "abc",
		PublishedStatus: event.Published,
	}
	err := eventService.UpdatePublishStatus(ctx, request)
	require.NoError(t, err)

	events, err := eventService.Events(ctx)
	evt := events[0]
	assert.Equal(t, int64(1), evt.ID)
	assert.Equal(t, event.Published, evt.PublishedStatus)
}
