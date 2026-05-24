package events

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

type EventFilters struct {
	Limit     int
	Type      string
	TraceID   string
	SessionID string
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Add(event Event) (Event, error) {
	metadataBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		return Event{}, err
	}

	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UnixMilli()
	}

	event.ReceivedAt = time.Now().UnixMilli()

	_, err = s.db.Exec(`
	INSERT INTO events (
		id,
		session_id,
		trace_id,
		name,
		timestamp,
		received_at,
		type,
		duration,
		metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		event.ID,
		event.SessionID,
		event.TraceID,
		event.Name,
		event.Timestamp,
		event.ReceivedAt,
		event.Type,
		event.Duration,
		string(metadataBytes),
	)

	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (s *Store) Clear() error {
	_, err := s.db.Exec(`DELETE FROM events`)
	return err
}

func (s *Store) GetAll(filters EventFilters) ([]Event, error) {
	limit := filters.Limit
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT id, session_id, trace_id, name, timestamp, received_at, type, duration, metadata
		FROM events
		WHERE 1=1
	`

	args := []any{}

	if filters.Type != "" {
		query += " AND type = ?"
		args = append(args, filters.Type)
	}

	if filters.TraceID != "" {
		query += " AND trace_id = ?"
		args = append(args, filters.TraceID)
	}

	if filters.SessionID != "" {
		query += " AND session_id = ?"
		args = append(args, filters.SessionID)
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var event Event
		var duration sql.NullFloat64
		var metadataString sql.NullString

		if err := rows.Scan(
			&event.ID,
			&event.SessionID,
			&event.TraceID,
			&event.Name,
			&event.Timestamp,
			&event.ReceivedAt,
			&event.Type,
			&duration,
			&metadataString,
		); err != nil {
			return nil, err
		}

		if duration.Valid {
			event.Duration = &duration.Float64
		}

		if metadataString.Valid && metadataString.String != "" {
			if err := json.Unmarshal([]byte(metadataString.String), &event.Metadata); err != nil {
				event.Metadata = map[string]interface{}{
					"parseError": err.Error(),
				}
			}
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}