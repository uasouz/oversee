package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"oversee/collector"
	"syscall"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	status "google.golang.org/grpc/status"
)

// Agent represents an agent in the system.
// This agent is responsible for receiving audit logs requests and processing them.
// The processing involves receiving the audit logs requests, parsing them, and then saving them to a local buffer.
// This buffer must be flushed periodically to ensure that the data is not lost.
// The buffer is also flushed when the agent is stopped for wathever reason.

type DispatchMode int

const (
	DispatchModeBatch = iota + 1
	DisptachModeIndividual
)

type Agent struct {
	Name         string
	Application  Application
	DispatchMode DispatchMode
	db           *badger.DB
	stream       *badger.Stream
	bufferFile   *os.File

	flushTime int

	signals chan os.Signal

	collectorClient     collector.CollectorClient
	collectorClientConn *grpc.ClientConn
	collectorAddress    string
}

type Application struct {
	Name          string
	Version       string //?? Maybe
	InitializedAt time.Time
}

func (agent *Agent) newCollectorClient() (collector.CollectorClient, error) {
	fmt.Println("Connecting to", agent.collectorAddress)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(agent.collectorAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	agent.collectorClientConn = conn

	c := collector.NewCollectorClient(conn)

	return c, nil
}

func (agent *Agent) shutdown() {
	// Flush and gracefully close the buffer
	agent.collectorClientConn.Close()

	agent.flushBuffer()
	agent.db.Close()

	// Exit with code 0
	os.Exit(0)
}

func (agent *Agent) Log(prefix string, log LogLine) error {
	return agent.db.Update(func(txn *badger.Txn) error {
		newLogId := uuid.New()
		err := txn.Set([]byte(prefix+"-"+newLogId.String()), log.Bytes())
		return err
	})
}

func (agent *Agent) batchDispatch(ctx context.Context, kvList *badger.KVList) error {
	logs := []*collector.PersistLogRequest{}

	for _, item := range kvList.GetKv() {
		logs = append(logs, &collector.PersistLogRequest{
			Id: string(item.Key),
		})
	}

	reply, err := agent.collectorClient.BatchPersistLog(ctx, &collector.BatchPersistLogRequest{
		Logs: logs,
	})

	if err != nil {
		return err
	}

	for _, result := range reply.Results {
		if result.GetSuccess() {
			fmt.Println("Persisted", result.Id)
			err = agent.db.Update(func(txn *badger.Txn) error {
				return txn.Delete([]byte(result.Id))
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (agent *Agent) simpleDispatch(ctx context.Context, kvList *badger.KVList) error {
	for _, item := range kvList.GetKv() {
		fmt.Println("Consuming", string(item.Key), string(item.Value))
		fmt.Println(agent.collectorAddress)

		reply, err := agent.collectorClient.PersistLog(ctx, &collector.PersistLogRequest{
			Id: string(item.Key),
		})

		if err != nil {
			st, ok := status.FromError(err)

			if ok {
				switch st.Code() {
				case codes.AlreadyExists:
					fmt.Println("Persisted", string(item.Key))
					err = agent.db.Update(func(txn *badger.Txn) error {
						return txn.Delete(item.Key)
					})

					if err != nil {
						return err
					}

					return nil
				default:
					return err
				}
			}

			return err
		}

		fmt.Println("Reply", reply)
		if reply.GetSuccess() && reply.Id == string(item.Key) {
			fmt.Println("Persisted", string(item.Key))
			err = agent.db.Update(func(txn *badger.Txn) error {
				return txn.Delete(item.Key)
			})

			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (agent *Agent) Start() error {
	signal.Notify(agent.signals, syscall.SIGTERM)
	defer agent.shutdown()

	var err error

	agent.collectorClient, err = agent.newCollectorClient()

	if err != nil {
		return err
	}

	// Start a timer to flush the buffer every minute
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Println("flushing..")
			agent.flushBuffer()
		}
	}()

	// Check if a file already exists and if it needs to be flushed
	// Create file if it does not exists
	db, err := badger.Open(badger.DefaultOptions("/tmp/trail"))

	if err != nil {
		return err
	}

	fmt.Println("Database OK")

	agent.db = db

	fmt.Println("Starting Stream")
	stream := db.NewStream()
	fmt.Println("Stream Started")

	stream.Send = func(buf *z.Buffer) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		kvList, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}

		if agent.DispatchMode == DispatchModeBatch {
			return agent.batchDispatch(ctx, kvList)
		}

		return agent.simpleDispatch(ctx, kvList)
	}

	agent.stream = stream

	<-agent.signals

	return nil
}

func (agent *Agent) flushBuffer() error {
	// Read registers on file and send to remote gRPC API to save Audit Logs
	// Open the buffer file for reading
	if err := agent.stream.Orchestrate(context.Background()); err != nil {
		return err
	}

	return nil
}

func NewAgent(name string, applicationName string, collectorAddress string) *Agent {
	agent := &Agent{
		Name:             name,
		signals:          make(chan os.Signal, 1),
		collectorAddress: collectorAddress,
		DispatchMode:     DispatchModeBatch,
		Application: Application{
			Name:          applicationName,
			Version:       "1.0.0",
			InitializedAt: time.Now(),
		},
	}

	return agent
}
