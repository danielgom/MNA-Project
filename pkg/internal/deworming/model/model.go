package model

import "time"

type Deworming struct {
	ID        int64      `db:"id"`
	PetID     int64      `db:"pet_id"`
	Address   string     `db:"address"`
	VetName   string     `db:"vet_name"`
	Date      *time.Time `db:"date"`
	NextDate  *time.Time `db:"next_date"`
	UpdatedAt *time.Time `db:"updated_at"`
}
