package persistence

import (
	"context"
	"oversee/core"
)

type LogPersistenceResult struct {
	ID      string
	Success bool
}

type Persistence interface {
	PersistLog(ctx context.Context, log *core.Log) (*LogPersistenceResult, error)
	BatchPersistLog(ctx context.Context, log []*core.Log) ([]*LogPersistenceResult, error)
}

