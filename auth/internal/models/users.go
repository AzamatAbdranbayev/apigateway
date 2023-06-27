package models

import (
	"math/rand"
	"time"
)

type User struct {
	Id          string    `json:"id"`
	FirstName   string    `json:"first_name" validate:"required"`
	LastName    string    `json:"last_name" validate:"required"`
	MiddleName  string    `json:"middle_name" validate:"required"`
	Login       string    `json:"login"`
	GroupNumber string    `json:"group_number" validate:"required"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserChangeBalanceRequest struct {
	Id      string  `json:"id" validate:"required"`
	Balance float64 `json:"balance" validate:"required"`
}

func (u *User) CreateLogin() {
	length := 8
	buf := make([]byte, length)
	digits := "0123456789"
	strAll := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = strAll[rand.Intn(len(strAll))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	u.Login = string(buf)
}
