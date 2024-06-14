package services

import (
	"context"
	"database/sql"

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

	_, err = db.Exec("INSERT INTO user_table (nama_user, password_user,email_user) VALUES (?, ?, ?)", log.Namauser, encpyt, log.Emailuser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})

	defer db.Close()

}

//JWT FUNC

// ---------------------------------------------------------------------------------------------------------------------
func GetUser(c *gin.Context) {
	// 1. Extract user data from request body
	var log models.UserData

	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// 2. Connect to database (consider using a connection pool)
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	defer db.Close() // Ensure database connection is closed

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(20 * time.Minute)

	// 3. Prepare SQL statement with named parameters (prevents SQL injection)
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT  email_user, password_user  FROM user_table")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Terjadi Kesalahan Saat Request DB"})
		return
	}

	defer rows.Close() // Ensure rows are closed

	// 4. Process Results
	found := false
	for rows.Next() {

		var email_user string
		var password_user string

		err := rows.Scan(&email_user, &password_user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
			return
		}

		// Compare password
		err = bcrypt.CompareHashAndPassword([]byte(password_user), []byte(log.Passuser))
		if err == nil {
			// If credentials match
			found = true
			tokenString, err := CreateToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"Message": "LOGIN BERHASIL", "TOKEN": tokenString})
			return
		}
	}

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Belum Terdaftar"})
	}
}

// ---------------------------------------------------------------------------------------------------------------------
