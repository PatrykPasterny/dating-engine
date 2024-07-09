package api

import (
	"context"

	pb "github.com/PatrykPasterny/dating-engine/tests/definition"
)

func (s *apiTestSuite) TestSuccessfullyCountLikedYou() {
	client := pb.NewExploreServiceClient(s.GrpcClient)

	request := pb.CountLikedYouRequest{
		RecipientUserId: "ab30308e-de0f-47df-9b51-55b9af86213d",
	}

	response, err := client.CountLikedYou(context.Background(), &request)
	if err != nil {
		s.T().Fatalf("failed counting users that liked the user: %v", err)
	}

	s.Equal(response.GetCount(), uint64(s.expectedUserLiked))
}
