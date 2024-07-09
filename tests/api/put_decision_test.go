package api

import (
	"context"

	pb "github.com/PatrykPasterny/dating-engine/tests/definition"
)

func (s *apiTestSuite) TestSuccessfullyPutDecisionOnDislikedUser() {
	client := pb.NewExploreServiceClient(s.GrpcClient)

	putRequest := pb.PutDecisionRequest{
		ActorUserId:     s.notLikedUserID,
		RecipientUserId: s.userID,
		LikedRecipient:  true,
	}

	putResponse, err := client.PutDecision(context.Background(), &putRequest)
	if err != nil {
		s.T().Fatalf("failed putting decision on user that did not like the user: %v", err)
	}

	actorSideMatch, err := s.getMatch(context.Background(), putRequest.ActorUserId, putRequest.RecipientUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as an actor: %v", err)
	}

	recipientSideMatch, err := s.getMatch(context.Background(), putRequest.RecipientUserId, putRequest.ActorUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as a recipient: %v", err)
	}

	s.Equal(actorSideMatch.Liked, true)
	s.Equal(actorSideMatch.Matched, true)
	s.Equal(recipientSideMatch.Liked, true)
	s.Equal(recipientSideMatch.Matched, true)
	s.Equal(putResponse.MutualLikes, true)

	putRequest.LikedRecipient = false

	putResponse, err = client.PutDecision(context.Background(), &putRequest)
	if err != nil {
		s.T().Fatalf("failed putting decision on user that did not like the user: %v", err)
	}

	actorSideMatch, err = s.getMatch(context.Background(), putRequest.ActorUserId, putRequest.RecipientUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as an actor: %v", err)
	}

	recipientSideMatch, err = s.getMatch(context.Background(), putRequest.RecipientUserId, putRequest.ActorUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as a recipient: %v", err)
	}

	s.Equal(actorSideMatch.Liked, false)
	s.Equal(actorSideMatch.Matched, false)
	s.Equal(recipientSideMatch.Liked, true)
	s.Equal(recipientSideMatch.Matched, false)
	s.Equal(putResponse.MutualLikes, false)
}

func (s *apiTestSuite) TestSuccessfullyPutDecisionOnDislikedRecipient() {
	client := pb.NewExploreServiceClient(s.GrpcClient)

	putRequest := pb.PutDecisionRequest{
		ActorUserId:     s.userID,
		RecipientUserId: s.notLikedByUserID,
		LikedRecipient:  true,
	}

	putResponse, err := client.PutDecision(context.Background(), &putRequest)
	if err != nil {
		s.T().Fatalf("failed putting decision on user that liked the user: %v", err)
	}

	actorSideMatch, err := s.getMatch(context.Background(), putRequest.ActorUserId, putRequest.RecipientUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as an actor: %v", err)
	}

	recipientSideMatch, err := s.getMatch(context.Background(), putRequest.RecipientUserId, putRequest.ActorUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as a recipient: %v", err)
	}

	s.Equal(actorSideMatch.Liked, true)
	s.Equal(actorSideMatch.Matched, true)
	s.Equal(recipientSideMatch.Liked, true)
	s.Equal(recipientSideMatch.Matched, true)
	s.Equal(putResponse.MutualLikes, true)

	putRequest.LikedRecipient = false

	putResponse, err = client.PutDecision(context.Background(), &putRequest)
	if err != nil {
		s.T().Fatalf("failed putting decision on user that liked the user: %v", err)
	}

	actorSideMatch, err = s.getMatch(context.Background(), putRequest.ActorUserId, putRequest.RecipientUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as an actor: %v", err)
	}

	recipientSideMatch, err = s.getMatch(context.Background(), putRequest.RecipientUserId, putRequest.ActorUserId)
	if err != nil {
		s.T().Fatalf("failed getting match for user as a recipient: %v", err)
	}

	s.Equal(actorSideMatch.Matched, false)
	s.Equal(actorSideMatch.Liked, false)
	s.Equal(recipientSideMatch.Matched, false)
	s.Equal(recipientSideMatch.Liked, true)
	s.Equal(putResponse.MutualLikes, false)
}
