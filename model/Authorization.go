package model

type AuthHeader struct {
	Authorization string `header: "Authorization" binding:"required"`
}
