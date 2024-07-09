package api

import (
	"context"

	"github.com/PatrykPasterny/dating-engine/internal/model"
)

type MatchRepository interface {
	GetLikedUser(ctx context.Context, userID, paginationToken string, limit int64) ([]model.Match, error)
	GetNewLikedUser(ctx context.Context, userID, paginationToken string, limit int64) ([]model.Match, error)
	CountLikedUser(ctx context.Context, userID string) (uint64, error)
	MakeDecision(ctx context.Context, userID, recipientID string, decision bool) (bool, error)
}
