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

func (m *SnippetModel) Insert(name, description, expired string) (int, error) {
	stmt := "insert into Film (name,description,time,count) VALUES ($1, $2, $3, $4) RETURNING id"
	intExpires, err := strconv.Atoi(expired)

	if err != nil {
		return 0, err
	}
	var lastIndex int

	err = m.DB.QueryRow(context.Background(), stmt, name, description, time.Now().AddDate(0, 0, intExpires), 10).Scan(&lastIndex)
	if err != nil {
		return 0, err
	}
	return int(lastIndex), nil
}

func (m *SnippetModel) Get(id int) (*models.Films, error) {
	stmt := "Select id, name, description, time, count FROM Film WHERE id = $1"

	row := m.DB.QueryRow(context.Background(), stmt, id)

	s := &models.Films{}

	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.Time, &s.Count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Films, error) {
	stmt := "Select id, name, description, time, count FROM Film WHERE time > CLOCK_TIMESTAMP() and count > 0 ORDER BY time DESC"

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Films{}

	for rows.Next() {
		s := &models.Films{}

		err = rows.Scan(&s.ID, &s.Name, &s.Description, &s.Time, &s.Count)
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

func (m *SnippetModel) Update(id int) string {

	check, err := m.Get(id)

	if err != nil {
		return "No films"
	} else {
		if check.Count != 0 {
			stmt := "Update Film set count = $1 where id = $2 RETURNING id"

			row := m.DB.QueryRow(context.Background(), stmt, check.Count-1, id)

			var id int

			err = row.Scan(&id)

			if err == nil {
				return "You successfully bought ticket"
			} else {
				return "Interval server error"
			}
		} else {
			return "No tickets"
		}
	}
}
