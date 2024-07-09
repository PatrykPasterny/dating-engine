package api

import (
	"context"
	"log/slog"

	pb "github.com/PatrykPasterny/dating-engine/transfer/protobuf/definition"

	"github.com/google/uuid"
)

func (es *ExploreServer) ListLikedYou(
	ctx context.Context,
	request *pb.ListLikedYouRequest,
) (*pb.ListLikedYouResponse, error) {
	loggerWithFields := es.logger.With(
		slog.String("recipient_id", request.RecipientUserId),
	)

	loggerWithFields.Info("retrieving list of all users that liked the user")

	paginationToken := uuid.Nil.String()

	if request.PaginationToken != nil {
		paginationToken = *request.PaginationToken
	}

	likedYouList, err := es.matchRepository.GetLikedUser(ctx, request.RecipientUserId, paginationToken, es.pageSize)
	if err != nil {
		loggerWithFields.Error("failed to get all users that liked the user", err)

		return nil, err
	}

	var response pb.ListLikedYouResponse

	response.Likers = make([]*pb.ListLikedYouResponse_Liker, 0, len(likedYouList))

	for i := range likedYouList {
		if i == len(likedYouList)-1 {
			response.NextPaginationToken = &likedYouList[i].ActorUserID
		}

		liker := &pb.ListLikedYouResponse_Liker{
			ActorId: likedYouList[i].ActorUserID,
		}

		response.Likers = append(response.Likers, liker)
	}

	loggerWithFields.Info("successfully retrieved list of users that liked the user")

	return &response, nil
}

func (es *ExploreServer) ListNewLikedYou(
	ctx context.Context,
	request *pb.ListLikedYouRequest,
) (*pb.ListLikedYouResponse, error) {
	loggerWithFields := es.logger.With(
		slog.String("recipient_id", request.RecipientUserId),
	)

	loggerWithFields.Info("retrieving list of new users that liked the user")

	paginationToken := uuid.Nil.String()

	if request.PaginationToken != nil {
		paginationToken = *request.PaginationToken
	}

	likedYouList, err := es.matchRepository.GetNewLikedUser(ctx, request.RecipientUserId, paginationToken, es.pageSize)
	if err != nil {
		loggerWithFields.Error("failed to get new users that liked the user", err)

		return nil, err
	}

	var response pb.ListLikedYouResponse

	response.Likers = make([]*pb.ListLikedYouResponse_Liker, 0, len(likedYouList))

	for i := range likedYouList {
		if i == len(likedYouList)-1 {
			response.NextPaginationToken = &likedYouList[i].ActorUserID
		}

		liker := &pb.ListLikedYouResponse_Liker{
			ActorId: likedYouList[i].ActorUserID,
		}

		response.Likers = append(response.Likers, liker)
	}

	loggerWithFields.Info("successfully retrieved list of new users that liked the user")

	return &response, nil
}

func (es *ExploreServer) CountLikedYou(
	ctx context.Context,
	request *pb.CountLikedYouRequest,
) (*pb.CountLikedYouResponse, error) {
	loggerWithFields := es.logger.With(
		slog.String("recipient_id", request.RecipientUserId),
	)

	loggerWithFields.Info("counting users that liked the user")

	count, err := es.matchRepository.CountLikedUser(ctx, request.RecipientUserId)
	if err != nil {
		loggerWithFields.Error("failed to count users that liked the user", err)

		return nil, err
	}

	response := pb.CountLikedYouResponse{
		Count: count,
	}

	loggerWithFields.Info("successfully counted users that liked the user")

	return &response, err
}

func (es *ExploreServer) PutDecision(
	ctx context.Context,
	request *pb.PutDecisionRequest,
) (*pb.PutDecisionResponse, error) {
	loggerWithFields := es.logger.With(
		slog.String("actor_id", request.ActorUserId),
		slog.String("recipient_id", request.RecipientUserId),
	)

	loggerWithFields.Info("applying new decision of the user")

	mutualLikes, err := es.matchRepository.MakeDecision(
		ctx,
		request.ActorUserId,
		request.RecipientUserId,
		request.LikedRecipient,
	)
	if err != nil {
		loggerWithFields.Error("failed to make decision on user", err)

		return nil, err
	}

	response := pb.PutDecisionResponse{
		MutualLikes: mutualLikes,
	}

	loggerWithFields.Info("successfully made new decision of user")

	return &response, err
}
