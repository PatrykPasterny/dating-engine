package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	"github.com/PatrykPasterny/dating-engine/internal/config"
	pb "github.com/PatrykPasterny/dating-engine/transfer/protobuf/definition"
)

type ExploreServer struct {
	pb.UnimplementedExploreServiceServer
	logger          *slog.Logger
	grpcServer      *grpc.Server
	matchRepository MatchRepository
	pageSize        int64
	baseURL         string
}

func NewExploreServer(
	logger *slog.Logger,
	cfg *config.Config,
	grpcServer *grpc.Server,
	repository MatchRepository,
	pageSize int64,
) *ExploreServer {
	return &ExploreServer{
		logger:          logger,
		grpcServer:      grpcServer,
		matchRepository: repository,
		baseURL:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		pageSize:        pageSize,
	}
}

func (es *ExploreServer) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signalChan

		es.logger.Info("received signal", slog.String("signal", sig.String()))
		cancel()
	}()

	wg.Add(1)

	go func(running *sync.WaitGroup) {
		defer running.Done()

		lis, err := net.Listen("tcp", fmt.Sprintf(es.baseURL))
		if err != nil {
			es.logger.Error(
				"failed to listen",
				err,
				slog.String("baseURL", es.baseURL),
			)

			cancel()
		}

		if err = es.grpcServer.Serve(lis); err != nil {
			es.logger.Error("failed to serve grpc", err)

			cancel()
		}

		es.logger.Info("explore service up and running")
	}(&wg)

	<-ctx.Done()

	es.grpcServer.Stop()

	wg.Wait()
}
