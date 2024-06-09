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

	_, err = db.Exec("INSERT INTO user_table (nama_user, password_user) VALUES (?, ?)", log.Namauser, encpyt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})

	defer db.Close()

}

//JWT FUNC

func CreateToken() {

}

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
	defer db.Close() // Ensure database connection is closed even on errors

	// 3. Prepare SQL statement with named parameters (prevents SQL injection)
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT password_user FROM user_table WHERE nama_user = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	defer stmt.Close() // Close prepared statement after use

	// 4. Execute query to get the stored hashed password
	var storedHashedPassword string
	err = stmt.QueryRowContext(ctx, log.Namauser).Scan(&storedHashedPassword)
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credential"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	// 5. Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(log.Passuser))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	// 6. Successful login
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
