package models

import "github.com/golang-jwt/jwt/v5"

type UserData struct {
	Namauser  string `json:"nama_user"`
	Passuser  string `json:"password_user"`
	Emailuser string `json:"email_user"`
}

type ClaimsData struct {
	NamaUser string `json:"nama_user"`
	jwt.RegisteredClaims
}
