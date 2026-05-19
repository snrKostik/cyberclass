package models

type Player struct {
	ID int64

	Nickname string
	RealName *string

	CreatedAt int64
}
