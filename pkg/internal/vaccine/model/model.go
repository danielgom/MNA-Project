package model

import "time"

type Vaccine struct {
	ID        int64      `db:"id"`
	PetID     int64      `db:"pet_id"`
	Type      string     `db:"type"`
	VetName   string     `db:"vet_name"`
	Address   string     `db:"address"`
	Date      *time.Time `db:"date"`
	NextDate  *time.Time `db:"next_date"`
	UpdatedAt *time.Time `db:"updated_at"`
}
