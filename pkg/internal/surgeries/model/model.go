package model

import "time"

type Surgery struct {
	ID        int64      `db:"id"`
	PetID     int64      `db:"pet_id"`
	Address   string     `db:"address"`
	VetName   string     `db:"vet_name"`
	Name      string     `db:"name"`
	Comments  string     `db:"comments"`
	Date      *time.Time `db:"date"`
	UpdatedAt *time.Time `db:"updated_at"`
}
