package model

type UserPayload struct {
	Id       int
	Iat      uint
	Name     string
	Email    string
	Phone    string
	Pst      string //staff position id
	Lst      string //staff location id
	Stp      string //staff type
	Scl      string //school id
	Uid      string //user id
	Uag      string //user agent
	Ip       string //Ip address
	Dacy     string //Default academic year
	AccessId int
}
