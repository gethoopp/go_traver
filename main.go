package main

import (
	"fmt"

	"github.com/gethoopp/go_traver.git/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	req := services.RegisterUserData
	res := services.GetUser

	r.POST("api/Register", req)
	r.POST("api/login", res)

	r.Run(":8080")

	fmt.Println()

}

//curl -X POST http://localhost:8080/api/login -d '{"email_user" : "Fisika@gmail.com","password_user" : "Fisika"}' -H "Content-Type: application/json"
