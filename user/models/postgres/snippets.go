package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"se09.com/user/models"
)

type Usermodel struct {
	DB *pgxpool.Pool
}

func (m *Usermodel) Insert(login, pass string) int {
	stmt := "INSERT INTO kino_user (login, pass, Balans)" + "VALUES($1, $2, $3) RETURNING id"

	var id int
	result := m.DB.QueryRow(context.Background(), stmt, login, pass, 0).Scan(&id)

	if result != nil {
		return 0
	}

	return id
}

func (m *Usermodel) Check(login, pass string) (*models.Kino_user, error) {
	stmt := "SELECT id, login, pass, Balans from  kino_user where login = $1 and pass = $2"
	s := &models.Kino_user{}
	err := m.DB.QueryRow(context.Background(), stmt, login, pass).Scan(&s.ID, &s.Login, &s.Pass, &s.Balans)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
