package models

//User model
type User struct {
	Addresses   []string `bson:"addresses" json:"addresses" binding:"required"`
	MainAddress string   `bson:"mainAddress" json:"mainAddress" binding:"required"`
	Username    string   `bson:"username" json:"username" binding:"required"`
	Email       string   `bson:"email" json:"email"`
	Phone       string   `bson:"phone" json:"phone"`
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
