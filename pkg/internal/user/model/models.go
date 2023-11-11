// Package model contains all models which are going to be inserted into the DB.
package model

import "time"

// User is the user which is going to be saved into the DB.
type User struct {
	ID        int64      `db:"id"`
	Password  string     `db:"password"`
	Name      string     `db:"name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	LastLogin *time.Time `db:"last_login"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
