package api

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	CustomerId string `json:"customerId" binding:"required"`
	Token      string `json:"token" binding:"required"`
}

type RegisterRequest struct {
	Customer
}
