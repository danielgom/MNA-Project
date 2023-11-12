package pet

import (
	"encoding/json"
	"strings"
	"time"
)

const timeParse = "2006-01-02"

type BirthDate time.Time

type RegisterRequest struct {
	Name      string    `json:"name"`
	Age       int8      `json:"age"`
	Breed     string    `json:"breed"`
	BirthDate BirthDate `json:"birth_date"`
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
	BirthDate    BirthDate  `json:"birth_date"`
	RegisterDate *time.Time `json:"register_date"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type Pets struct {
	Pets []*GeneralResponse `json:"pets"`
}

func (d *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(timeParse, s)
	if err != nil {
		return err
	}
	*d = BirthDate(t)
	return nil
}

func (d *BirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d).Format(timeParse))
}

func (d *BirthDate) Format(s string) string {
	t := time.Time(*d)
	return t.Format(s)
}
