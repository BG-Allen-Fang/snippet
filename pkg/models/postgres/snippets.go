package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"se09.com/pkg/models"
	"strconv"
	"time"
)

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires)" + "VALUES($1, $2, $3, $4) RETURNING id"
	date, err := strconv.Atoi(expires)

	if err != nil {
		return 0, err
	}

	var id int
	result := m.DB.QueryRow(context.Background(), stmt, title, content, time.Now(), time.Now().AddDate(0, 0, date)).Scan(&id)

	if result != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > CLOCK_TIMESTAMP() AND id = $1"
	s := &models.Snippet{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > CLOCK_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
