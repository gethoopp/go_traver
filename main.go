package main

import (
	"fmt"

	"github.com/gethoopp/go_traver.git/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	req := services.RegisterUserData
	res := services.LoginUser

	r.POST("api/Register", req)
	r.GET("api/login", res)

	r.Run(":8080")

	fmt.Println()

}

//curl -X POST http://localhost:8080/api/Register -d '{"nama_user":"Farahdiba", "password_user":"Farahdiba21042001"}' -H "Content-Type: application/json"
