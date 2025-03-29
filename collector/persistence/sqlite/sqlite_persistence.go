package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"oversee/collector/persistence"
	"oversee/core"

	_ "github.com/mattn/go-sqlite3"
)

type SQLitePersistence struct {
	db *sql.DB
}

func NewSQLitePersistence(dbPath string) (*SQLitePersistence, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = createLogTable(db); err != nil {
		return nil, fmt.Errorf("failed to create log table: %w", err)
	}

	return &SQLitePersistence{db: db}, nil
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (s *SQLitePersistence) PersistLog(ctx context.Context, log *core.Log) (*persistence.LogPersistenceResult, error) {
	query := `
		INSERT INTO logs (
			id,
			timestamp,
			service_name,
			operation,
			actor_id,
			actor_type,
			affected_resources,
			metadata,
			integrity_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	timestampInt := log.Timestamp.Unix()
	metadataJSON, err := json.Marshal(log.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	affectedResourcesJSON, err := json.Marshal(log.AffectedResources)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal affected resources: %w", err)
	}

	res, err := stmt.ExecContext(ctx,
		log.ID,
		timestampInt,
		log.ServiceName,
		log.Operation,
		log.ActorID,
		log.ActorType,
		affectedResourcesJSON,
		string(metadataJSON),
		log.IntegrityHash,
	)
	if err != nil {

		if isUniqueConstraintError(err) {
			return nil, core.ErrorAlreadyPersistedLog
		}

		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &persistence.LogPersistenceResult{
		ID:      fmt.Sprintf("%d", id),
		Success: true,
	}, nil
}

func (s *SQLitePersistence) BatchPersistLog(ctx context.Context, logs []*core.Log) ([]*persistence.LogPersistenceResult, error) {
	query := `
		INSERT INTO logs (
			id,
			timestamp,
			service_name,
			operation,
			actor_id,
			actor_type,
			affected_resources,
			metadata,
			integrity_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	results := make([]*persistence.LogPersistenceResult, len(logs))
	for i, log := range logs {
		timestampInt := log.Timestamp.Unix()
		metadataJSON, err := json.Marshal(log.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}

		affectedResourcesJSON, err := json.Marshal(log.AffectedResources)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal affected resources: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			log.ID,
			timestampInt,
			log.ServiceName,
			log.Operation,
			log.ActorID,
			log.ActorType,
			affectedResourcesJSON,
			string(metadataJSON),
			log.IntegrityHash,
		)
		if err != nil {
			results[i] = &persistence.LogPersistenceResult{
				ID:      log.ID,
				Success: false,
			}
		} else {
			results[i] = &persistence.LogPersistenceResult{
				ID:      log.ID,
				Success: true,
			}
		}
	}

	return results, nil
}

func createLogTable(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS logs (
	id CHAR(36) PRIMARY KEY,
	timestamp INTEGER NOT NULL,
	service_name TEXT NOT NULL,
	operation TEXT NOT NULL,
	actor_id TEXT NOT NULL,
	actor_type TEXT NOT NULL,
	affected_resources TEXT NOT NULL,
	metadata TEXT NOT NULL,
	integrity_hash TEXT NOT NULL
)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create log table: %w", err)
	}
	return nil
}
