package agent

import (
	"context"
	"fmt"
	"net"
	"oversee/core"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type IngestionAPI struct {
	UnimplementedAgentServer
	agent *Agent
}

// Log implements AgentServer.
func (a IngestionAPI) Log(ctx context.Context, request *LogRequest) (*LogReply, error) {
	err := a.agent.Log(&core.Log{
		ID:                uuid.New(),
		Timestamp:         request.Timestamp.AsTime(),
		ServiceName:       request.ServiceName,
		Operation:         request.Operation,
		ActorId:           request.ActorId,
		ActorType:         request.ActorType,
		AffectedResources: []string{},
		Metadata:          map[string]any{},
		IntegrityHash:     "",
	})

	return &LogReply{
		Success: err == nil,
	}, err
}

func (a *IngestionAPI) Serve() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 4092))

	if err != nil {
		return err
	}

	s := grpc.NewServer()

	RegisterAgentServer(s, a)
	go func() {
		err = s.Serve(listener)
		if err != nil {
			fmt.Println(err)
		}
	}()

	a.agent.Start()

	return nil
}

func NewIngestionAPI(collectorAddress string) *IngestionAPI {
	agent := NewAgent("main", "demo", collectorAddress)

	return &IngestionAPI{
		agent: agent,
	}
}
