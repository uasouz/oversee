package core

import "time"

type Log struct {
	ID                string
	Timestamp         time.Time
	ServiceName       string
	Operation         string
	ActorID           string
	ActorType         string
	AffectedResources []string
	Metadata          map[string]any
	IntegrityHash     string
}
