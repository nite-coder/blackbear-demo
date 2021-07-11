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
	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/internal/pkg/database"
	internalDatabase "github.com/nite-coder/blackbear-demo/internal/pkg/database"

	"github.com/nite-coder/blackbear-demo/pkg/domain"
	"github.com/nite-coder/blackbear-demo/pkg/event/proto"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	eventDatabase "github.com/nite-coder/blackbear-demo/pkg/event/repository/mysql"
	eventUsecase "github.com/nite-coder/blackbear-demo/pkg/event/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener

	_db  *gorm.DB
	_cfg config.Configuration
	// repo
	_eventRepo domain.EventRepository

	// services
	_eventService domain.EventUsecase

	// grpc server
	_eventServer eventProto.EventServiceServer

	// grpc client
	_eventClient eventProto.EventServiceClient
)

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(m *testing.M) {
	var err error
	_cfg = config.New("app_test.yml")

	_cfg.InitLogger("event")

	// initial database
	_db, err = database.InitDatabase(_cfg, "starter_db")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// repo
	_eventRepo := eventDatabase.NewEventRepository(_cfg, _db)

	// services
	_eventService := eventUsecase.NewEventUsecase(_cfg, _db, _eventRepo)

	// grpc server
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    (time.Duration(5) * time.Second), // Ping the client if it is idle for 5 seconds to ensure the connection is still active
				Timeout: (time.Duration(5) * time.Second), // Wait 5 second for the ping ack before assuming the connection is dead
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             (time.Duration(2) * time.Second), // If a client pings more than once every 2 seconds, terminate the connection
				PermitWithoutStream: true,                             // Allow pings even when there are no active streams
			},
		),
		grpc.ChainUnaryInterceptor(
			Interceptor(),
		),
	)
	_eventServer := NewEventServer(_cfg, _eventService)
	proto.RegisterEventServiceServer(grpcServer, _eventServer)

	go func() {
		lis = bufconn.Listen(bufSize)
		if err := grpcServer.Serve(lis); err != nil {
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

	grpcServer.GracefulStop()

	os.Exit(exitVal)
}

func TestGetEvents(t *testing.T) {
	// clear database data
	err := internalDatabase.RunSQLScripts(_db, _cfg.Path("test", "database", "starter_db"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := _eventClient.GetEvents(ctx, &eventProto.GetEventsRequest{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(resp.Events))

	evt := resp.Events[0]
	assert.Equal(t, int64(1), evt.Id)
	assert.Equal(t, proto.PublishedStatus_PublishedStatus_Draft, evt.PublishedStatus)
}

func TestUpdatePublishStatus(t *testing.T) {
	// clear database data
	err := internalDatabase.RunSQLScripts(_db, _cfg.Path("test", "database", "starter_db"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	request := proto.UpdatePublishStatusRequest{
		EventId:         1,
		TransId:         "abc",
		PublishedStatus: proto.PublishedStatus_PublishedStatus_Published,
	}

	_, err = _eventClient.UpdatePublishStatus(ctx, &request)
	require.NoError(t, err)

	getEventResp, err := _eventClient.GetEvents(ctx, &eventProto.GetEventsRequest{})
	require.NoError(t, err)

	assert.Equal(t, 1, len(getEventResp.Events))

	evt := getEventResp.Events[0]
	assert.Equal(t, int64(1), evt.Id)
	assert.Equal(t, proto.PublishedStatus_PublishedStatus_Published, evt.PublishedStatus)

	_, err = _eventClient.UpdatePublishStatus(ctx, &request)
	assert.EqualError(t, err, "rpc error: code = Internal desc = status of resource is wrong")
}
