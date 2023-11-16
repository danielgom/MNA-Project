package model

import "time"

type VetVisit struct {
	ID        int64      `db:"id"`
	PetID     int64      `db:"pet_id"`
	Address   string     `db:"address"`
	VetName   string     `db:"vet_name"`
	Reason    string     `db:"reason"`
	Comments  string     `db:"comments"`
	Date      *time.Time `db:"date"`
	UpdatedAt *time.Time `db:"updated_at"`
}
