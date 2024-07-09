package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/PatrykPasterny/dating-engine/tests/common"
	"github.com/PatrykPasterny/dating-engine/tests/model"
)

type apiTestSuite struct {
	*common.TestSuite
	userID                 string
	notLikedByUserID       string
	notLikedUserID         string
	expectedMutualLikes    int
	expectedLikedByUser    int
	expectedUserLiked      int
	expectedMutualDislikes int
}

func NewApiTestSuite() (*apiTestSuite, error) {
	const (
		testUserID         = "ab30308e-de0f-47df-9b51-55b9af86213d"
		testMutualLikes    = 20
		testLikedByUser    = 40
		testUserLiked      = 30
		testMutualDislikes = 20
	)

	ts, err := common.NewTestSuite()
	if err != nil {
		return nil, fmt.Errorf("creating new test suite: %w", err)
	}

	testSuite := &apiTestSuite{
		TestSuite:              ts,
		userID:                 testUserID,
		expectedMutualLikes:    testMutualLikes,
		expectedLikedByUser:    testLikedByUser,
		expectedUserLiked:      testUserLiked,
		expectedMutualDislikes: testMutualDislikes,
	}

	return testSuite, nil
}

func TestApi(t *testing.T) {
	s, err := NewApiTestSuite()
	if err != nil {
		t.Fatalf("unable to create test suite: %v", err)

		return
	}

	suite.Run(t, s)
}

func (s *apiTestSuite) SetupTest() {
	if err := s.initializeDatabase(); err != nil {
		s.FailNow("failed initializing database", err)
	}
}

func (s *apiTestSuite) TearDownTest() {
	s.notLikedUserID = ""
	s.notLikedByUserID = ""

	filter := bson.D{
		{
			Key: "actorUserID", Value: bson.D{
				{
					Key:   "$gt",
					Value: uuid.Nil.String(),
				},
			},
		},
	}

	if _, err := s.Collection.DeleteMany(context.Background(), filter, options.Delete()); err != nil {
		s.FailNow("unable to delete all matches from collection", err)
	}
}

func (s *apiTestSuite) initializeDatabase() error {
	matches := make(
		[]interface{},
		0,
		s.expectedMutualLikes+s.expectedLikedByUser+s.expectedUserLiked+s.expectedMutualDislikes,
	)

	// generate users that recipient user liked and they liked back
	for range s.expectedMutualLikes {
		actorUserID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating new actorUserID: %w", err)
		}

		userLike := bson.D{
			{
				Key:   "recipientUserID",
				Value: s.userID,
			},
			{
				Key:   "actorUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "liked",
				Value: true,
			},
			{
				Key:   "matched",
				Value: true,
			},
		}

		actorLike := bson.D{
			{
				Key:   "recipientUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "actorUserID",
				Value: s.userID,
			},
			{
				Key:   "liked",
				Value: true,
			},
			{
				Key:   "matched",
				Value: true,
			},
		}

		matches = append(matches, userLike, actorLike)
	}

	// generate users that recipient user did not like and they did not like him back
	for range s.expectedMutualDislikes {
		actorUserID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating new actorUserID: %w", err)
		}

		userDislike := bson.D{
			{
				Key:   "recipientUserID",
				Value: s.userID,
			},
			{
				Key:   "actorUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "liked",
				Value: false,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		actorDislike := bson.D{
			{
				Key:   "recipientUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "actorUserID",
				Value: s.userID,
			},
			{
				Key:   "liked",
				Value: false,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		matches = append(matches, userDislike, actorDislike)
	}

	// generate users that recipient user did not like or see yet, but they did like him
	for i := range s.expectedUserLiked - s.expectedMutualLikes {
		actorUserID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating new actorUserID: %w", err)
		}

		userLike := bson.D{
			{
				Key:   "recipientUserID",
				Value: s.userID,
			},
			{
				Key:   "actorUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "liked",
				Value: true,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		actorDislike := bson.D{
			{
				Key:   "recipientUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "actorUserID",
				Value: s.userID,
			},
			{
				Key:   "liked",
				Value: false,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		matches = append(matches, userLike)

		// We want to simulate situation where user disliked(made decision on) only half of the users
		// that liked him back
		if i%2 == 0 {
			if s.notLikedByUserID == "" {
				s.notLikedByUserID = actorUserID.String()
			}

			matches = append(matches, actorDislike)
		}
	}

	// generate users that recipient user liked, but they did not like or see yet
	for i := range s.expectedLikedByUser - s.expectedMutualLikes {
		actorUserID, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating new actorUserID: %w", err)
		}

		userDislike := bson.D{
			{
				Key:   "recipientUserID",
				Value: s.userID,
			},
			{
				Key:   "actorUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "liked",
				Value: false,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		actorLike := bson.D{
			{
				Key:   "recipientUserID",
				Value: actorUserID.String(),
			},
			{
				Key:   "actorUserID",
				Value: s.userID,
			},
			{
				Key:   "liked",
				Value: true,
			},
			{
				Key:   "matched",
				Value: false,
			},
		}

		matches = append(matches, actorLike)

		// We want to simulate situation where only half of the actors disliked(made decision on) the test user
		if i%2 == 0 {
			if s.notLikedUserID == "" {
				s.notLikedUserID = actorUserID.String()
			}

			matches = append(matches, userDislike)
		}
	}

	if _, err := s.Collection.InsertMany(context.Background(), matches, options.InsertMany()); err != nil {
		return fmt.Errorf("inserting many matches into collections: %w", err)
	}

	return nil
}

func (s *apiTestSuite) getMatch(ctx context.Context, actorID, recipientID string) (*model.Match, error) {
	var match model.Match

	filter := bson.D{
		{
			Key:   "actorUserID",
			Value: actorID,
		},
		{
			Key:   "recipientUserID",
			Value: recipientID,
		},
	}

	result := s.Collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, fmt.Errorf("finding match: %w", result.Err())
	}

	err := result.Decode(&match)
	if err != nil {
		return nil, fmt.Errorf("deciding match: %w", err)
	}

	return &match, nil
}
