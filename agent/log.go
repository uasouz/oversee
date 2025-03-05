package agent

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LogLine struct {
	ID              uuid.UUID
	Timestamp         time.Time
	ServiceName       string
	Operation         string
	ActorId           string
	ActorType         string
	AffectedResources []string
	Metadata          map[string]any
	IntegrityHash     string
}

func (l *LogLine) String() string {
	s, _ := json.Marshal(l)
	return string(s)
}

func (l *LogLine) Bytes() []byte {
	s, _ := json.Marshal(l)
	return s
}
