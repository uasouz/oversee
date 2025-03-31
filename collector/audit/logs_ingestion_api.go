package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"oversee/collector/persistence"
	"oversee/core"

	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogsIngestionAPI struct {
	UnimplementedCollectorServer
	persistence persistence.Persistence
}

// BatchPersistLog implements CollectorServer.
func (c LogsIngestionAPI) BatchPersistLog(ctx context.Context, request *BatchPersistLogRequest) (*PersistLogsReply, error) {
	fmt.Printf("Persisting %d logs", len(request.Logs))
	logs := []*core.Log{}

	for _, log := range request.Logs {
		logs = append(logs, LogEntityFromAPILog(log))
	}

	results, err := c.persistence.BatchPersistLog(ctx, logs)

	if err != nil {
		return nil, err
	}

	replies := []*PersistLogReply{}

	for _, result := range results {
		replies = append(replies, LogPersistenceResultToPersistLogReply(result))
	}

	v, _ := json.Marshal(results)
	fmt.Println(replies, string(v))

	return &PersistLogsReply{
		Results: replies,
	}, nil
}

func LogPersistenceResultToPersistLogReply(result *persistence.LogPersistenceResult) *PersistLogReply {
	reply := &PersistLogReply{
		Id:      result.ID,
		Success: result.Success,
	}

	if result.Reason != nil {
		reply.Reason = &Error{
			Message: result.Reason.Message,
			Code:    int32(result.Reason.Code),
		}
	}

	return reply
}

func LogEntityFromAPILog(log *Log) *core.Log {
	metadataMap := make(map[string]any)
	if log.Metadata != nil {
		for key, value := range log.Metadata.Fields {
			metadataMap[key] = value.AsInterface()
		}
	}

	logUUID, _ := uuid.Parse(log.Id)

	return &core.Log{
		ID:                logUUID,
		Timestamp:         log.Timestamp.AsTime(),
		ServiceName:       log.ServiceName,
		Operation:         log.Operation,
		ActorId:           log.ActorId,
		ActorType:         log.ActorType,
		AffectedResources: log.AffectedResources,
		Metadata:          metadataMap,
		IntegrityHash:     log.IntegrityHash,
	}
}

// PersistLog implements CollectorServer.
func (c LogsIngestionAPI) PersistLog(ctx context.Context, request *PersistLogRequest) (*PersistLogReply, error) {
	fmt.Println("Persisting", request.Log.Id)

	if request.Log == nil {
		return nil, fmt.Errorf("Log required")
	}

	log := LogEntityFromAPILog(request.Log)

	result, err := c.persistence.PersistLog(ctx, log)

	if err != nil {
		if err == core.ErrorAlreadyPersistedLog {
			return &PersistLogReply{
				Id:      request.Log.Id,
				Success: true,
			}, status.Error(codes.AlreadyExists, err.Error())
		}
		return &PersistLogReply{
			Id:      request.Log.Id,
			Success: false,
		}, err
	}

	return LogPersistenceResultToPersistLogReply(result), nil
}

func (a *LogsIngestionAPI) Serve() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 4093))
	fmt.Println("Starting collector API on port 4093")

	if err != nil {
		return err
	}

	s := grpc.NewServer()

	RegisterCollectorServer(s, a)

	err = s.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func NewLogsIngestionAPI(persistence persistence.Persistence) *LogsIngestionAPI {
	return &LogsIngestionAPI{
		persistence: persistence,
	}
}
