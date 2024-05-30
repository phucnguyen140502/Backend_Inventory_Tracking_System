package controllers

import (
	"backend/database"
	"backend/models"
	"backend/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignUp(c *gin.Context) {
	var body struct {
		Name     string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := security.HashPassword(body.Password)

	user := models.Users{
		FullName:  body.Name,
		Email:     body.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := database.DB.Create(&user); result.Error != nil {

		if result.Error.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exists",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Users
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !security.ComparePasswords(user.Password, body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong password",
		})
		return
	}

	// Tao JWT Authentication Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"name":     user.FullName,
		"email":    user.Email,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(security.SECRET_KEY))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create token",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func GetAllUsers(c *gin.Context) {
	var users []models.Users
	database.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
