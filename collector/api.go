package collector

import (
	"context"
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
func (c CollectorAPI) BatchPersistLog(context.Context, *BatchPersistLogRequest) (*PersistLogReply, error) {
	panic("unimplemented")
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

	return &PersistLogReply{
		Id:      result.ID,
		Success: true,
	}, nil

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
