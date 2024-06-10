package services

import (
	"context"
	"database/sql"
	"fmt"

	"net/http"
	"time"

	"github.com/gethoopp/go_traver.git/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserData(c *gin.Context) {
	// Connect to the database
	var log models.UserData
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		panic(err)

	}
	// Close the database after the function

	// Set database connection settings
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(20 * time.Minute)

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	//bcrypt atau enkripsi password dari user
	encpyt, err := bcrypt.GenerateFromPassword([]byte(log.Passuser), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}
	// Insert data

	_, err = db.Exec("INSERT INTO user_table (nama_user, password_user) VALUES (?, ?)", log.Namauser, encpyt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})

	defer db.Close()

}

//JWT FUNC

// ---------------------------------------------------------------------------------------------------------------------
func LoginUser(c *gin.Context) {
	// 1. Extract user data from request body
	var log models.UserData

	// 2. Connect to database (consider using a connection pool)
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(20 * time.Minute)

	// 3. Prepare SQL statement with named parameters (prevents SQL injection)
	ctx := context.Background()
	stmt, err := db.QueryContext(ctx, "SELECT iduser_table, nama_user , password_user FROM user_table")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Terjadi Kesalahan Saat Request DB"})
		return
	}
	defer stmt.Close() // Close prepared statement after use

	fmt.Println(log.Namauser)

	var result []gin.H

	for stmt.Next() {
		var id int
		var nama_user string
		var password_user string

		err := stmt.Scan(&id, &nama_user, &password_user)

		if err != nil {
			panic(err)
		}

		result := append(result, gin.H{"Namauser": nama_user, "id": id})

		password := "Farahdiba21042001"

		err = bcrypt.CompareHashAndPassword([]byte(password_user), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, result)

	}

	// 5. Compare the stored hashed password with the provided password

	// 6. Successful login

	defer db.Close() // Ensure database connection is closed even on errors
}
