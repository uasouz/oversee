package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID                uuid.UUID
	Timestamp         time.Time
	ServiceName       string
	Operation         string
	ActorId           string
	ActorType         string
	AffectedResources []string
	Metadata          map[string]any
	IntegrityHash     string
}

func (l *Log) String() string {
	return string(l.Bytes())
}

func (l *Log) Bytes() []byte {
	s, _ := json.Marshal(l)
	return s
}
