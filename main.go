package main

import (
	"context"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/PatrykPasterny/dating-engine/internal/config"
	"github.com/PatrykPasterny/dating-engine/internal/repository"
	"github.com/PatrykPasterny/dating-engine/transfer/protobuf/api"
	pb "github.com/PatrykPasterny/dating-engine/transfer/protobuf/definition"
)

const configPath = "./internal/config/config.yml"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	logger.Info("the explorer service is starting")

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		logger.Error("failed getting configuration", err)

		return
	}

	clientOpts := options.Client().ApplyURI(cfg.Database.URI)

	mongoClient, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		logger.Error("failed connecting to mongoDB instance", err)

		return
	}

	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("failed disconnecting from mongoDB instance", err)

			return
		}
	}()

	collection := mongoClient.Database(cfg.Database.Name).Collection(cfg.Database.Collection)

	exploreRepository := repository.NewExploreRepository(mongoClient, collection)

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	exploreServer := api.NewExploreServer(logger, cfg, grpcServer, exploreRepository, cfg.PageSize)
	pb.RegisterExploreServiceServer(grpcServer, exploreServer)

	exploreServer.Run()
}
