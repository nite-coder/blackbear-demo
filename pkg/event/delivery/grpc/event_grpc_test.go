package grpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event/proto"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	cfg := config.Configuration{}
	eventService := eventService.NewEventService(cfg)

	s := grpc.NewServer()
	eventServer := NewEventServer(cfg, eventService)
	proto.RegisterEventServiceServer(s, eventServer)

	go func() {
		lis = bufconn.Listen(bufSize)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func TestGetEvents(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewEventServiceClient(conn)
	resp, err := client.GetEvents(ctx, &empty.Empty{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(resp.Data))

	evt := resp.Data[0]
	assert.Equal(t, int64(1), evt.Id)
	assert.Equal(t, proto.PublishedStatus_Draft, evt.PublishedStatus)
}

func TestUpdatePublishStatus(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewEventServiceClient(conn)

	request := proto.UpdatePublishStatusRequest{
		EventId:         1,
		TransId:         "abc",
		PublishedStatus: proto.PublishedStatus_Published,
	}

	_, err = client.UpdatePublishStatus(ctx, &request)
	require.NoError(t, err)

	getEventResp, err := client.GetEvents(ctx, &empty.Empty{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(getEventResp.Data))

	evt := getEventResp.Data[0]
	assert.Equal(t, int64(1), evt.Id)
	assert.Equal(t, proto.PublishedStatus_Published, evt.PublishedStatus)
}
