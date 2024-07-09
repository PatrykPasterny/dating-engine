package common

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestSuite struct {
	suite.Suite
	dbClient   *mongo.Client
	Collection *mongo.Collection
	Logger     *slog.Logger
	GrpcClient *grpc.ClientConn
}

func NewTestSuite() (*TestSuite, error) {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	cfg := NewConfig()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(cfg.BaseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating new grpc dbClient: %w", err)
	}

	clientOpts := options.Client().ApplyURI(cfg.DatabaseURI)

	dbClient, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongoDB instance: %w", err)
	}

	collection := dbClient.Database(cfg.DatabaseName).Collection(cfg.DatabaseCollection)

	ts := &TestSuite{
		dbClient:   dbClient,
		Collection: collection,
		Logger:     logger,
		GrpcClient: conn,
	}

	return ts, nil
}

func (ts *TestSuite) SetupSuite() {
	ts.Logger.Info("setting up the test suite")
}

func (ts *TestSuite) TearDownSuite() {
	ts.Logger.Info("tearing down the test suite")

	if err := ts.dbClient.Disconnect(context.Background()); err != nil {
		ts.Logger.Error("failed disconnecting from mongoDB instance", err)

		return
	}

	if err := ts.GrpcClient.Close(); err != nil {
		ts.Logger.Error("failed closing grpc client instance", err)

		return
	}
}
