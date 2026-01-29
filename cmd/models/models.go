package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	PWD       string    `json:"pwd"`
	Role      int       `json:"role"`
	Account   int       `json:"account"`
	Picture   string    `json:"picture"`
	Token     string    `json:"token"`
	Status    int       `json:"status"`
	Confirmed int       `json:"confirmed"`
	Code      string    `json:"code"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Active    int       `json:"active"`
}
type Feedback struct {
	Id      int       `json:"id"`
	IdUser  string    `json:"id_user"`
	Name    string    `json:"name"`
	Picture string    `json:"picture"`
	Rate    int       `json:"stars"`
	Comment string    `json:"feedback"`
	Consent int       `json:"consent"`
	Display int       `json:"display"`
	Watched int       `json:"watched"`
	Created time.Time `json:"created"`
}
