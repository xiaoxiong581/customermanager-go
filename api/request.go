package api

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	CustomerId string `json:"customerId"`
	token      string `json:"token"`
}

type RegisterRequest struct {
	Customer
}
