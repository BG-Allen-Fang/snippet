package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Films struct {
	ID          int
	Name        string
	Description string
	Time        time.Time
	Count       int
}

type Ticket struct {
	ID    int
	U_id  int
	Name  string
	Time  string
	Price int
}

type Kino_user struct {
	ID     int
	Login  string
	Pass   string
	Balans int
}
