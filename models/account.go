package models

type Account struct {
	ID int64

	Username     string
	PasswordHash string

	CreatedAt int64
}
