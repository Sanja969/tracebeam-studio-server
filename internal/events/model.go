package events

type Event struct {
	ID        string                 `json:"id"`
	TraceID   string                 `json:"traceId,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
	Name      string                 `json:"name"`
	Timestamp int64                  `json:"timestamp"`
	Type      string                 `json:"type"`
	Duration  *float64              `json:"duration,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ReceivedAt int64 `json:"receivedAt"`
}