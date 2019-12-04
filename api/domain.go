package api

type Customer struct {
	CustomerName string `json:"customerName" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
}
