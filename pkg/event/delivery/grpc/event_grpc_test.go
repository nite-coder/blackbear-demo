package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/event"
	"github.com/jasonsoft/starter/pkg/event/proto"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	eventDatabase "github.com/jasonsoft/starter/pkg/event/repository/database"
	eventService "github.com/jasonsoft/starter/pkg/event/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
	// repo
	_eventRepo event.Repository

	// services
	_eventService event.Servicer

	// grpc server
	_eventServer eventProto.EventServiceServer

	// grpc client
	_eventClient eventProto.EventServiceClient
)

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(m *testing.M) {
	cfg := config.New("app.yml")

	cfg.InitLogger("event")

	// initial database
	db, err := cfg.InitDatabase("starter")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// repo
	_eventRepo := eventDatabase.NewEventRepository(cfg, db)

	// services
	_eventService := eventService.NewEventService(cfg, _eventRepo)

	// grpc server
	s := grpc.NewServer()
	_eventServer := NewEventServer(cfg, _eventService)
	proto.RegisterEventServiceServer(s, _eventServer)

	go func() {
		lis = bufconn.Listen(bufSize)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// grpc client
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_eventClient = proto.NewEventServiceClient(conn)

	exitVal := m.Run()

	s.GracefulStop()

	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func TestGetEvents(t *testing.T) {
	ctx := context.Background()

	resp, err := _eventClient.GetEvents(ctx, &eventProto.GetEventsRequest{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(resp.Events))

	evt := resp.Events[0]
	assert.Equal(t, int64(1), evt.Id)
	assert.Equal(t, proto.PublishedStatus_PublishedStatus_Draft, evt.PublishedStatus)
}

// func TestUpdatePublishStatus(t *testing.T) {
// 	ctx := context.Background()

// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
// 	if err != nil {
// 		t.Fatalf("Failed to dial bufnet: %v", err)
// 	}
// 	defer conn.Close()
// 	client := proto.NewEventServiceClient(conn)

// 	request := proto.UpdatePublishStatusRequest{
// 		EventId:         1,
// 		TransId:         "abc",
// 		PublishedStatus: proto.PublishedStatus_Published,
// 	}

// 	_, err = client.UpdatePublishStatus(ctx, &request)
// 	require.NoError(t, err)

// 	getEventResp, err := client.GetEvents(ctx, &empty.Empty{})
// 	require.NoError(t, err)

// 	assert.Equal(t, 1, len(getEventResp.Data))

// 	evt := getEventResp.Data[0]
// 	assert.Equal(t, int64(1), evt.Id)
// 	assert.Equal(t, proto.PublishedStatus_Published, evt.PublishedStatus)
// }
