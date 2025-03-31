package persistence

import (
	"context"
	"oversee/core"
)

type LogPersistenceResult struct {
	ID      string
	Success bool
	Reason  *core.Error
}

// TODO: Add support to supabase persistence

type SearchQuery struct {
	ServiceName       string
	Operation         string
	ActorID           string
	ActorType         string
	AffectedResources []string
	Metadata          map[string]any
	CursorTimestamp   int64
	CursorID          string
}

type Persistence interface {
	PersistLog(ctx context.Context, log *core.Log) (*LogPersistenceResult, error)
	BatchPersistLog(ctx context.Context, log []*core.Log) ([]*LogPersistenceResult, error)
	ListLogs(ctx context.Context, cursorTimestamp int64, cursorID string) ([]*core.Log, error)
	SearchLogs(ctx context.Context, query SearchQuery) ([]*core.Log, error)
}
