package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"oversee/collector/persistence"
	"oversee/core"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CollectorAPI struct {
	UnimplementedCollectorServer
	persistence persistence.Persistence
}

// BatchPersistLog implements CollectorServer.
func (c CollectorAPI) BatchPersistLog(ctx context.Context, request *BatchPersistLogRequest) (*PersistLogsReply, error) {
	fmt.Printf("Persisting %d logs", len(request.Logs))
	logs := []*core.Log{}
	for _, persistLogRequest := range request.Logs {
		logs = append(logs, LogFromPersistLogRequest(persistLogRequest))
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
	return &PersistLogReply{
		Id:      result.ID,
		Success: result.Success,
	}
}

func LogFromPersistLogRequest(request *PersistLogRequest) *core.Log {
	metadataMap := make(map[string]any)
	if request.Metadata != nil {
		for key, value := range request.Metadata.Fields {
			metadataMap[key] = value.AsInterface()
		}
	}

	return &core.Log{
		ID:                request.Id,
		Timestamp:         request.Timestamp.AsTime(),
		ServiceName:       request.ServiceName,
		Operation:         request.Operation,
		ActorID:           request.ActorId,
		ActorType:         request.ActorType,
		AffectedResources: request.AffectedResources,
		Metadata:          metadataMap,
		IntegrityHash:     request.IntegrityHash,
	}
}

// PersistLog implements CollectorServer.
func (c CollectorAPI) PersistLog(ctx context.Context, request *PersistLogRequest) (*PersistLogReply, error) {
	fmt.Println("Persisting", request.Id)

	log := LogFromPersistLogRequest(request)

	result, err := c.persistence.PersistLog(ctx, log)

	if err != nil {
		if err == core.ErrorAlreadyPersistedLog {
			return &PersistLogReply{
				Id:      request.Id,
				Success: true,
			}, status.Error(codes.AlreadyExists, err.Error())
		}
		return &PersistLogReply{
			Id:      request.Id,
			Success: false,
		}, err
	}

	return LogPersistenceResultToPersistLogReply(result), nil
}

func (a *CollectorAPI) Serve() error {
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

func NewCollectorAPI(persistence persistence.Persistence) *CollectorAPI {
	return &CollectorAPI{
		persistence: persistence,
	}
}
