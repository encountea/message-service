package repository

import (
	"time"

	"github.com/encountea/message-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveMessage(msg models.Message) error {
	_, err := r.db.Exec("INSERT INTO messages (content, processed, created_at) VALUES ($1, $2, $3)",
		msg.Content, msg.Processed, time.Now())
	return err
}

func (r *Repository) MarkAsProcessed(id int) error {
	_, err := r.db.Exec("UPDATE messages SET processed = TRUE WHERE id = $1", id)
	return err
}

func (r *Repository) GetProcessedCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM messages WHERE processed = TRUE").Scan(&count)
	return count, err
}
