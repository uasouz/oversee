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

func (s *SQLitePersistence) ListLogs(ctx context.Context) ([]*core.Log, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, timestamp, service_name, operation, actor_id, actor_type, affected_resources, metadata, integrity_hash FROM logs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*core.Log
	for rows.Next() {
		log := &core.Log{}
		var metadataJSON string
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.ServiceName, &log.Operation, &log.ActorId, &log.ActorType, &log.AffectedResources, &metadataJSON, &log.IntegrityHash); err != nil {
			return nil, err
		}
		log.Metadata = map[string]any{}
		if err := json.Unmarshal([]byte(metadataJSON), &log.Metadata); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *SQLitePersistence) SearchLogs(ctx context.Context, query persistence.SearchQuery) ([]*core.Log, error) {
	var whereClauses []string
	var args []any

	if query.ServiceName != "" {
		whereClauses = append(whereClauses, "service_name = ?")
		args = append(args, query.ServiceName)
	}

	if query.Operation != "" {
		whereClauses = append(whereClauses, "operation = ?")
		args = append(args, query.Operation)
	}

	if query.ActorID != "" {
		whereClauses = append(whereClauses, "actor_id = ?")
		args = append(args, query.ActorID)
	}

	if query.ActorType != "" {
		whereClauses = append(whereClauses, "actor_type = ?")
		args = append(args, query.ActorType)
	}

	if len(query.AffectedResources) > 0 {
		resourcesJSON, err := json.Marshal(query.AffectedResources)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal affected resources: %w", err)
		}
		whereClauses = append(whereClauses, "affected_resources = ?")
		args = append(args, string(resourcesJSON))
	}

	if len(query.Metadata) > 0 {
		metadataJSON, err := json.Marshal(query.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		whereClauses = append(whereClauses, "metadata LIKE ?")
		args = append(args, "%"+string(metadataJSON)+"%")
	}

	queryString := "SELECT id, timestamp, service_name, operation, actor_id, actor_type, affected_resources, metadata, integrity_hash FROM logs"
	if len(whereClauses) > 0 {
		queryString += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := s.db.QueryContext(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*core.Log
	for rows.Next() {
		log := &core.Log{}
		var metadataJSON string
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.ServiceName, &log.Operation, &log.ActorId, &log.ActorType, &log.AffectedResources, &metadataJSON, &log.IntegrityHash); err != nil {
			return nil, err
		}
		log.Metadata = map[string]any{}
		if err := json.Unmarshal([]byte(metadataJSON), &log.Metadata); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
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
		log.ActorId,
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

func (s *SQLitePersistence) areLogsAlreadyPersisted(ctx context.Context, logIDs []string) (map[string]bool, error) {
	existingLogs := make(map[string]bool)
	for _, id := range logIDs {
		var exists bool
		err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM logs WHERE id = ?)", id).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check if log with ID %s exists: %w", id, err)
		}
		existingLogs[id] = exists
	}
	return existingLogs, nil
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
			log.ActorId,
			log.ActorType,
			affectedResourcesJSON,
			string(metadataJSON),
			log.IntegrityHash,
		)
		if err != nil {
			results[i] = &persistence.LogPersistenceResult{
				ID:      log.ID.String(),
				Success: false,
				Reason:  core.ErrorAlreadyPersistedLog,
			}
		} else {
			results[i] = &persistence.LogPersistenceResult{
				ID:      log.ID.String(),
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
