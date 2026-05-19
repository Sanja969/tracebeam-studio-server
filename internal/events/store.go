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
		trace_id,
		name,
		timestamp,
		received_at,
		type,
		duration,
		metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`,
		event.ID,
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

func (s *Store) GetAll() ([]Event, error) {
	rows, err := s.db.Query(`
		SELECT id, trace_id, name, timestamp, received_at, type, duration, metadata
		FROM events
		ORDER BY timestamp ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var event Event
		var metadataString sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.TraceID,
			&event.Name,
			&event.Timestamp,
			&event.ReceivedAt,
			&event.Type,
			&event.Duration,
			&metadataString,
		)
		if err != nil {
			return nil, err
		}

		if metadataString.Valid && metadataString.String != "" {
			json.Unmarshal([]byte(metadataString.String), &event.Metadata)
		}

		events = append(events, event)
	}

	return events, nil
}
