package api

type Customer struct {
	CustomerName string `json:"customerName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}