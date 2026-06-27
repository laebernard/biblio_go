package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	UserID  int
	IsAdmin bool
	jwt.RegisteredClaims
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	res, err := DB.Exec(
		"INSERT INTO users (name, email, password, isAdmin) VALUES (?, ?, ?, ?)",
		input.Name,
		input.Email,
		string(hashedPassword),
		false,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	userID, _ := res.LastInsertId()

	// JWT
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:  int(userID),
		IsAdmin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var id int
	var hashedPassword string
	var isAdmin bool

	err := DB.QueryRow(
		"SELECT id, password, isAdmin FROM users WHERE email = ?",
		input.Email,
	).Scan(&id, &hashedPassword, &isAdmin)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Générer JWT
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:  id,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

type UpdateMeInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UpdateMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UpdateMeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password si fourni
	var hashedPassword string
	if input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		hashedPassword = string(hash)
	}

	// Update dynamique
	query := "UPDATE users SET name = ?, email = ?"
	args := []interface{}{input.Name, input.Email}

	if hashedPassword != "" {
		query += ", password = ?"
		args = append(args, hashedPassword)
	}

	query += " WHERE id = ?"
	args = append(args, userID)

	_, err := DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
	})
}

func UpdateUserByAdmin(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	id := c.Param("id")

	var targetIsAdmin bool
	err := DB.QueryRow("SELECT isAdmin FROM users WHERE id = ?", id).Scan(&targetIsAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetIsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify an admin"})
		return
	}

	var input UpdateMeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var hashedPassword string
	if input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		hashedPassword = string(hash)
	}

	query := "UPDATE users SET name = ?, email = ?"
	args := []interface{}{input.Name, input.Email}

	if hashedPassword != "" {
		query += ", password = ?"
		args = append(args, hashedPassword)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err = DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated by admin",
	})
}

func ResetDatabase(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")

	if !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	err := ResetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reset failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database reset successful",
	})
}
