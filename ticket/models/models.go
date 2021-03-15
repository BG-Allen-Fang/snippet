package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Ticket struct {
	ID    int
	U_id  int
	Name  string
	Time  string
	Price int
}
