package security

import (
	"net/http"
	"time"

	"biblio_go/database"
	"biblio_go/logger"

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

// Register godoc
// @Summary      Inscription d'un nouvel utilisateur
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterInput  true  "Données d'inscription"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /user/register [post]
func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Register - Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	logger.Info("Register - New registration attempt for email: %s", input.Email)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	res, err := database.DB.Exec(
		"INSERT INTO users (name, email, password, isAdmin) VALUES (?, ?, ?, ?)",
		input.Name,
		input.Email,
		string(hashedPassword),
		false,
	)

	if err != nil {
		logger.Error("Register - Email already exists: %s", input.Email)
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

	logger.Info("Register - User registered successfully: %s (ID: %d)", input.Email, userID)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary      Connexion utilisateur
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginInput  true  "Email et mot de passe"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /user/login [post]
func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Login - Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	logger.Info("Login - Login attempt for email: %s", input.Email)

	var id int
	var hashedPassword string
	var isAdmin bool

	err := database.DB.QueryRow(
		"SELECT id, password, isAdmin FROM users WHERE email = ?",
		input.Email,
	).Scan(&id, &hashedPassword, &isAdmin)

	if err != nil {
		logger.Error("Login - Invalid credentials for email: %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		logger.Error("Login - Invalid password for email: %s", input.Email)
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

	logger.Info("Login - User logged in successfully: %s", input.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

type UpdateMeInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateMe godoc
// @Summary      Met à jour son propre profil
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      UpdateMeInput  true  "Nouvelles données"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur): Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /user/me [put]
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

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
	})
}

// UpdateUserByAdmin godoc
// @Summary      Met a jour un utilisateur (admin uniquement)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string         true  "User ID"
// @Param        input  body      UpdateMeInput  true  "Donnees utilisateur"
// @Param        Authorization  header    string  true  "Bearer token admin: Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{id} [put]
func UpdateUserByAdmin(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")
	if !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	id := c.Param("id")

	var targetIsAdmin bool
	err := database.DB.QueryRow("SELECT isAdmin FROM users WHERE id = ?", id).Scan(&targetIsAdmin)
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

	_, err = database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated by admin",
	})
}

// ResetDatabase godoc
// @Summary      Reinitialise la base de donnees (admin uniquement)
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token admin: Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /reset [delete]
func ResetDatabase(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")

	if !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	err := database.ResetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reset failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database reset successful",
	})
}
