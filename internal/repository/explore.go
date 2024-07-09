package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/PatrykPasterny/dating-engine/internal/model"
)

type ExploreRepository struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

func NewExploreRepository(mongoClient *mongo.Client, collection *mongo.Collection) *ExploreRepository {
	return &ExploreRepository{
		mongoClient: mongoClient,
		collection:  collection,
	}
}

func (er *ExploreRepository) GetLikedUser(
	ctx context.Context,
	userID, paginationToken string,
	limit int64,
) ([]model.Match, error) {
	findOptions := options.Find().SetSort(
		bson.D{
			{
				Key:   "actorUserID",
				Value: 1,
			},
		},
	).SetLimit(limit)

	fmt.Println(userID)

	filters := bson.D{
		{
			Key: "recipientUserID", Value: userID,
		},
		{
			Key: "liked", Value: true,
		},
		{
			Key: "actorUserID", Value: bson.D{
				{
					Key: "$gt", Value: paginationToken,
				},
			},
		},
	}

	cur, err := er.collection.Find(ctx, filters, findOptions)
	if err != nil {
		return nil, fmt.Errorf("finding users that liked the user: %w", err)
	}

	var likedUser []model.Match

	if err = cur.All(ctx, &likedUser); err != nil {
		return nil, fmt.Errorf("retrieving all users that liked the user: %w", err)
	}

	return likedUser, nil
}

func (er *ExploreRepository) GetNewLikedUser(
	ctx context.Context,
	userID, paginationToken string,
	limit int64,
) ([]model.Match, error) {
	findOptions := options.Find().SetSort(
		bson.D{
			{
				Key:   "actorUserID",
				Value: 1,
			},
		},
	).SetLimit(limit)

	filters := bson.D{
		{
			Key: "recipientUserID", Value: userID,
		},
		{
			Key: "liked", Value: true,
		},
		{
			Key: "matched", Value: false,
		},
		{
			Key: "actorUserID", Value: bson.D{
				{
					Key: "$gt", Value: paginationToken,
				},
			},
		},
	}

	cur, err := er.collection.Find(ctx, filters, findOptions)
	if err != nil {
		return nil, fmt.Errorf("finding new users that liked the user: %w", err)
	}

	var newLikedUser []model.Match

	if err = cur.All(ctx, &newLikedUser); err != nil {
		return nil, fmt.Errorf("retrieving all new users that liked the user: %w", err)
	}

	return newLikedUser, nil
}

func (er *ExploreRepository) CountLikedUser(ctx context.Context, userID string) (uint64, error) {
	filters := bson.D{
		{
			Key: "recipientUserID", Value: userID,
		},
		{
			Key: "liked", Value: true,
		},
	}

	count, err := er.collection.CountDocuments(ctx, filters)
	if err != nil {
		return 0, fmt.Errorf("counting users that liked the user: %w", err)
	}

	return uint64(count), nil
}

func (er *ExploreRepository) MakeDecision(
	ctx context.Context,
	userID, recipientID string,
	decision bool,
) (bool, error) {
	findOptions := options.FindOne()

	userFilters := bson.D{
		{
			Key: "actorUserID", Value: userID,
		},
		{
			Key: "recipientUserID", Value: recipientID,
		},
	}

	recipientFilters := bson.D{
		{
			Key: "recipientUserID", Value: userID,
		},
		{
			Key: "actorUserID", Value: recipientID,
		},
	}

	session, err := er.mongoClient.StartSession()
	if err != nil {
		return false, fmt.Errorf("starting new mongo session: %w", err)
	}

	var mutualLikes bool

	if err = session.StartTransaction(); err != nil {
		return false, fmt.Errorf("starting new mongo transaction: %w", err)
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		userResult := er.collection.FindOne(ctx, userFilters, findOptions)
		if userResult.Err() != nil {
			return fmt.Errorf("finding user that made new decision: %w", userResult.Err())
		}

		recipientResult := er.collection.FindOne(ctx, recipientFilters, findOptions)
		if recipientResult.Err() != nil {
			return fmt.Errorf("finding user that recieved new decision: %w", recipientResult.Err())
		}

		var userMatch, recipientMatch model.Match

		if err = userResult.Decode(&userMatch); err != nil {
			return fmt.Errorf("decoding user that made new decision: %w", err)
		}

		if err = recipientResult.Decode(&recipientMatch); err != nil {
			return fmt.Errorf("decoding user that recieved new decision: %w", err)
		}

		mutualLikes = decision && recipientMatch.Liked

		updateUser := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "liked",
						Value: decision,
					},
					{
						Key:   "matched",
						Value: mutualLikes,
					},
				},
			},
		}

		updateRecipient := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{
						Key:   "matched",
						Value: mutualLikes,
					},
				},
			},
		}

		if _, err = er.collection.UpdateOne(ctx, userFilters, updateUser, options.Update()); err != nil {
			return fmt.Errorf("updating user with new decision: %w", err)
		}

		if _, err = er.collection.UpdateOne(ctx, recipientFilters, updateRecipient, options.Update()); err != nil {
			return fmt.Errorf("updating recipient with new decision: %w", err)
		}

		return nil
	}); err != nil {
		return false, fmt.Errorf("performing mongo transaction: %w", err)
	}
	session.EndSession(ctx)

	return mutualLikes, nil
}
