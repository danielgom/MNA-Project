package pet

import "time"

type RegisterRequest struct {
	Name      string     `json:"name"`
	Age       int8       `json:"age"`
	Breed     string     `json:"breed"`
	BirthDate *time.Time `json:"birth_date"`
}

type UpdateRequest struct {
	Name      string     `json:"name"`
	UserID    int64      `json:"-"`
	Age       int8       `json:"age"`
	Breed     string     `json:"breed"`
	BirthDate *time.Time `json:"birth_date"`
}

type GeneralResponse struct {
	ID           int64      `json:"id"`
	OwnerName    string     `json:"owner_name"`
	OwnerID      int64      `json:"owner_id"`
	Name         string     `json:"name"`
	Age          int8       `json:"age"`
	Breed        string     `json:"breed"`
	BirthDate    *time.Time `json:"birth_date"`
	RegisterDate *time.Time `json:"register_date"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type Pets struct {
	Pets []*GeneralResponse `json:"pets"`
}
