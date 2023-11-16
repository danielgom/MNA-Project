package surgeries

import (
	"encoding/json"
	"strings"
	"time"
)

const timeParse = "2006-01-02"

type CommonDate time.Time

type RegisterRequest struct {
	VetName  string      `json:"vet_name"`
	Address  string      `json:"address"`
	Name     string      `json:"name"`
	Comments string      `json:"comments"`
	Date     *CommonDate `json:"date"`
}

type UpdateRequest struct {
	UserID   int64       `json:"-"`
	PetID    int64       `json:"pet_id"`
	VetName  string      `json:"vet_name"`
	Address  string      `json:"address"`
	Name     string      `json:"name"`
	Comments string      `json:"comments"`
	Date     *CommonDate `json:"date"`
}

type GeneralResponse struct {
	ID       int64       `json:"id"`
	PetName  string      `json:"pet_name"`
	PetID    int64       `json:"pet_id"`
	VetName  string      `json:"vet_name"`
	Address  string      `json:"address"`
	Name     string      `json:"name"`
	Comments string      `json:"comments"`
	Date     *CommonDate `json:"date"`
}

type Surgeries struct {
	Surgeries []*GeneralResponse `json:"surgeries"`
}

func (d *CommonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(timeParse, s)
	if err != nil {
		return err
	}
	*d = CommonDate(t)
	return nil
}

func (d *CommonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d).Format(timeParse))
}

func (d *CommonDate) Format(s string) string {
	t := time.Time(*d)
	return t.Format(s)
}
