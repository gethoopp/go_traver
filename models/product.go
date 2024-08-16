package models

import "github.com/golang-jwt/jwt/v5"

type UserData struct {
	UserID    int    `json:"iduser_table"`
	Namauser  string `json:"nama_user"`
	Passuser  string `json:"password_user"`
	Emailuser string `json:"email_user"`
}

//struct data product

type Product struct {
	IDProduct   int    `json:"idproduct_table"`
	Name        string `json:"product_name"`
	DetailsProd string `json:"product_detail"`
}

// Struct data untuk JWT
type ClaimsData struct {
	NamaUser string `json:"nama_user"`
	jwt.RegisteredClaims
}

// Struct data untuk Respons Payment
type ResponsePayment struct {
	Code   int `json:"code"`
	Status int `json:"status"`
	Data   int `json:"data"`
}

// Struct data untuk err respons
type ResponseErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

//Struct data untuk req Payment

type ResponseReqPayment struct {
	UserId int `json:"user_id" binding:"required"`
	Amount int `json:"amount"  binding:"required"`
}

// Struct data untuk token Respons
type TokenResponse struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}
