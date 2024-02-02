package model

import "time"

type Entities struct {
	Id           string `db:"primary"`
	FullName     string
	ShortName    string
	ParentId     string
	EntityTypeId string
	Email        string
	Phone        string
	Country      string
	Operator     string
	Status       uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
