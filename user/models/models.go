package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Kino_user struct {
	ID     int
	Login  string
	Pass   string
	Balans int
}
