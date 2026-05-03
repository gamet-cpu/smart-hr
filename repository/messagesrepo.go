package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

func (m *MessageRepository) SentReq(ctx context.Context, name, email, topic, description string) error {
	_, err := m.db.Exec(ctx,
		`INSERT INTO messages (name, email, topic, description)
		 VALUES ($1,$2,$3,$4)`,
		name, email, topic, description)

	return err
}
