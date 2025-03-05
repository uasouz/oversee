package collector

import (
	"context"
	"fmt"
	"net"

	grpc "google.golang.org/grpc"
)

type CollectorAPI struct {
	UnimplementedCollectorServer
}

// BatchPersistLog implements CollectorServer.
func (c CollectorAPI) BatchPersistLog(context.Context, *BatchPersistLogRequest) (*PersistLogReply, error) {
	panic("unimplemented")
}

// PersistLog implements CollectorServer.
func (c CollectorAPI) PersistLog(ctx context.Context, request *PersistLogRequest) (*PersistLogReply, error) {
	fmt.Println("Persisting", request.Id)
	return &PersistLogReply{
		Id:      request.Id,
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

func NewCollectorAPI() *CollectorAPI {
	return &CollectorAPI{}
}
