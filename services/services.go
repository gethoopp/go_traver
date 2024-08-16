package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"net/http"
	"time"

	"github.com/gethoopp/go_traver.git/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserData(c *gin.Context) {
	var log models.UserData
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	encpyt, err := bcrypt.GenerateFromPassword([]byte(log.Passuser), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(20 * time.Minute)

	done := make(chan error)
	defer close(done)

	go func() {
		_, err = db.Exec("INSERT INTO user_table (nama_user, password_user,email_user) VALUES (?, ?, ?)", log.Namauser, encpyt, log.Emailuser)
		done <- err // megnirim data ke dalam channel done
	}()

	select {
	case err := <-done: //mengambil data dari channel done
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Request timed out"})
	}
}

// ---------------------------------------------------------------------------------------------------------------------
func GetUser(c *gin.Context) {
	var log models.UserData
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(20 * time.Minute)

	done := make(chan error)
	defer close(done)
	var email_user, password_user string

	go func() {
		rows, err := db.Query("SELECT email_user, password_user FROM user_table")
		if err != nil {
			done <- err
			return
		}
		defer rows.Close()

		found := false
		for rows.Next() {
			err := rows.Scan(&email_user, &password_user)
			if err != nil {
				done <- err
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(password_user), []byte(log.Passuser))
			if err == nil {
				found = true
				break
			}
		}

		if !found {
			done <- nil
		} else {
			tokenString, err := CreateToken()
			if err != nil {
				done <- err
			} else {
				done <- nil
				c.JSON(http.StatusOK, gin.H{"Message": "LOGIN BERHASIL", "TOKEN": tokenString})
			}
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error processing data"})
		} else if email_user == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Belum Terdaftar"})
		}
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Request timed out"})
	}
}

//Get all data from product

func GetData(c *gin.Context) {

	ctx := context.Background()

	// Membuka koneksi database
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/travel_db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	chanGet := make(chan error)

	var id int
	var name string
	var DetailsProd string
	defer close(chanGet)
	var result []gin.H

	go func() {
		rows, err := db.QueryContext(ctx, "SELECT idproduct_table, product_name, product_detail FROM product_table")
		if err != nil {
			chanGet <- err
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id, &name, &DetailsProd)
			if err != nil {
				chanGet <- err
				return
			}

			var DetailsProdMap map[string]interface{}

			if err := json.Unmarshal([]byte(DetailsProd), &DetailsProdMap); err != nil {
				chanGet <- err
				return
			}

			result = append(result, gin.H{"id": id, "NameProd": name, "DetailsProd": DetailsProdMap})
		}

		chanGet <- nil
	}()

	select {
	case err := <-chanGet:
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Message": "Tidak Berhasil Mendapatkan Data", "Error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Request timed out"})
	}
}

// Payment Gateway menggunakan midtrans


