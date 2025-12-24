package repository

import (
	"context"
	"time"

	"pastebin-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository struct {
	db *pgxpool.Pool
}

func NewPasteRepository(db *pgxpool.Pool) *PasteRepository {
	return &PasteRepository{db: db}
}

// create inserts new paste into database
func (r *PasteRepository) Create(ctx context.Context, paste *models.Paste) error {
	query := `
        INSERT INTO pastes (content, created_at) 
        VALUES ($1, $2) 
        RETURNING id, created_at
    `

	err := r.db.QueryRow(ctx, query, paste.Content, time.Now()).Scan(
		&paste.ID,
		&paste.CreatedAt,
	)

	return err
}

// getByID retrieves a paste by id
func (r *PasteRepository) GetByID(ctx context.Context, id int) (*models.Paste, error) {
	paste := &models.Paste{}
	query := `SELECT id, content, created_at FROM pastes WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&paste.ID,
		&paste.Content,
		&paste.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // not found
	}

	if err != nil {
		return nil, err
	}

	return paste, nil
}

// update modifies existing paste
func (r *PasteRepository) Update(ctx context.Context, paste *models.Paste) error {
	query := `UPDATE pastes SET content = $1 WHERE id = $2`

	cmdTag, err := r.db.Exec(ctx, query, paste.Content, paste.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// delete removes paste from database
func (r *PasteRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM pastes WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// getAll retrieves all pastes
func (r *PasteRepository) GetAll(ctx context.Context) ([]models.Paste, error) {
	query := `SELECT id, content, created_at FROM pastes ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []models.Paste
	for rows.Next() {
		var paste models.Paste
		if err := rows.Scan(&paste.ID, &paste.Content, &paste.CreatedAt); err != nil {
			return nil, err
		}
		pastes = append(pastes, paste)
	}

	return pastes, rows.Err()
}
