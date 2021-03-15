package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"se09.com/ticket/models"
	"time"
)

type TicketModel struct {
	DB *pgxpool.Pool
}

func (m *TicketModel) Insert(name, expired string, u_id, price int) (int, error) {
	stmt := "insert into Ticket (u_id,name,time,price) VALUES ($1, $2, $3, $4) RETURNING id"

	var lastIndex int
	expired = expired[:len(expired)-10]
	err := m.DB.QueryRow(context.Background(), stmt, u_id, name, expired, price).Scan(&lastIndex)
	if err != nil {
		return 0, err
	}
	return lastIndex, nil
}

func (m *TicketModel) Latest(id int) ([]*models.Ticket, error) {
	stmt := "Select id, u_id, name, time, price FROM Ticket WHERE u_id = $1"

	rows, err := m.DB.Query(context.Background(), stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Ticket{}

	for rows.Next() {
		s := &models.Ticket{}
		var t time.Time
		err = rows.Scan(&s.ID, &s.U_id, &s.Name, &t, &s.Price)
		s.Time = t.String()
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
