package events

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/mzeahmed/noticeal/internal/database/sqlc"
)

// Service contains the business logic of the events module. It consumes
// the sqlc-generated Queries directly, with no repository layer in
// between.
type Service struct {
	queries *sqlc.Queries
}

// NewService creates a new events service backed by db.
func NewService(db *sql.DB) *Service {
	return &Service{queries: sqlc.New(db)}
}

// Create persists e and returns it with its generated fields (ID,
// CreatedAt) populated.
func (s *Service) Create(ctx context.Context, e Event) (Event, error) {
	data, err := marshalData(e.Data)
	if err != nil {
		return Event{}, err
	}

	row, err := s.queries.CreateEvent(ctx, sqlc.CreateEventParams{
		Source:  e.Source,
		Type:    e.Type,
		Status:  e.Status,
		Title:   e.Title,
		Message: e.Message,
		Data:    data,
	})
	if err != nil {
		return Event{}, err
	}

	e.ID = row.ID
	e.CreatedAt = row.CreatedAt

	return e, nil
}

// marshalData encodes data as JSON for storage in the events.data column,
// leaving it NULL when there is nothing to store.
func marshalData(data map[string]string) (sql.NullString, error) {
	if len(data) == 0 {
		return sql.NullString{}, nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return sql.NullString{}, err
	}

	return sql.NullString{String: string(b), Valid: true}, nil
}
