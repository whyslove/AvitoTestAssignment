package models

type User struct {
	Id      int     `json:"-"`
	Balance float64 `json:"balance" binding:"required"`
}

type UserOperation struct {
	UserId        int     `json:"id" binding:"required"`
	AmountOfMoney float64 `json:"amount_money" binding:"required"`
}

type UserTransferOperation struct {
	SenderId      int     `json:"sender_id" binding:"required"`
	ReceiverId    int     `json:"receiver_id" binding:"required"`
	AmountOfMoney float64 `json:"amount_money" binding:"required"`
}
