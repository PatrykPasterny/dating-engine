package api

import (
	"context"

	pb "github.com/PatrykPasterny/dating-engine/tests/definition"
)

func (s *apiTestSuite) TestSuccessfullyGetLikedYou() {
	client := pb.NewExploreServiceClient(s.GrpcClient)

	request := pb.ListLikedYouRequest{
		RecipientUserId: "ab30308e-de0f-47df-9b51-55b9af86213d",
	}

	response, err := client.ListLikedYou(context.Background(), &request)
	if err != nil {
		s.T().Fatalf("failed getting list of users that liked the user: %v", err)
	}

	responseLength := len(response.Likers)

	// check whether pagination is also working
	for len(response.GetLikers()) > 0 {
		paginationToken := response.NextPaginationToken

		request = pb.ListLikedYouRequest{
			RecipientUserId: "ab30308e-de0f-47df-9b51-55b9af86213d",
			PaginationToken: paginationToken,
		}

		response, err = client.ListLikedYou(context.Background(), &request)
		if err != nil {
			s.T().Fatalf("failed getting list of users that liked the user: %v", err)
		}

		responseLength += len(response.GetLikers())
	}

	s.Equal(responseLength, s.expectedUserLiked)
}
