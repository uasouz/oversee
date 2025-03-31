package audit

import (
	"context"
	"oversee/collector/persistence"
	"oversee/core"
)

type SearchService struct {
	persistence persistence.Persistence
}

func NewSearchService(p persistence.Persistence) *SearchService {
	return &SearchService{
		persistence: p,
	}
}

func (s *SearchService) SearchLogs(ctx context.Context, query persistence.SearchQuery) ([]*core.Log, error) {
	return s.persistence.SearchLogs(ctx, query)
}

func (s *SearchService) ListLogs(ctx context.Context, cursorTimestamp int64, cursorID string) ([]*core.Log, error) {
	return s.persistence.ListLogs(ctx, cursorTimestamp, cursorID)
}
