package controller

/* import (
	"net/http"

	"github.com/gethoopp/go_traver.git/models"
	"github.com/gethoopp/go_traver.git/services"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var log models.UserData

	// Bind JSON payload to log variable
	if err := c.ShouldBindJSON(&log); err != nil {
		// Return 400 status code for bad request
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	// Register user data
	if err := services.RegisterUserData(log); err != nil {
		// Handle potential error from RegisterUserData
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register user"})
		return
	}

	// Return success response

}*/
