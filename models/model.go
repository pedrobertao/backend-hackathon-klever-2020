package models

import (
	"time"
)

type infoStatus int

// InfoStatus enum
const (
	Inactive infoStatus = iota + 1
	Pending
	Active
)

// UserPhone model
type UserPhone struct {
	Phone  string     `bson:"phone" json:"phone" binding:"min=5,max=20"`
	Status infoStatus `bson:"status" json:"status" binding:"oneof=1,2,3"`
}

// UserEmail model
type UserEmail struct {
	Email  string     `bson:"email" json:"email" binding:"email"`
	Status infoStatus `bson:"status" json:"status" binding:"oneof=1,2,3"`
}

//User model
type User struct {
	Addresses   []string  `bson:"addresses" json:"addresses" binding:"required,min=1"`
	MainAddress string    `bson:"mainAddress" json:"mainAddress" binding:"required"`
	Username    string    `bson:"username" json:"username" binding:"required"`
	Email       UserEmail `bson:"email" json:"email"`
	Phone       UserPhone `bson:"phone" json:"phone"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

//SMSRequest model
type SMSRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Coin   string  `json:"coin" binding:"required"`
}

//SMS model
type SMS struct {
	To   string `json:"to" binding:"required"`
	From string `json:"from" binding:"required"`
	Body string `json:"body" binding:"required"`
}
