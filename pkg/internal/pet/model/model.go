package model

import "time"

type Pet struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	Name         string     `db:"name"`
	Age          int8       `db:"age"`
	Breed        string     `db:"breed"`
	Type         string     `db:"type"`
	Sex          string     `db:"sex"`
	Color        string     `db:"color"`
	BirthDate    *time.Time `db:"birth_date"`
	RegisterDate *time.Time `db:"register_date"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
