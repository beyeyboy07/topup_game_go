package controllers

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"topup_games_go/models"
	"topup_games_go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Name     string `json:"name" example:"Budi"`
	Email    string `json:"email" example:"budi@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"budi@example.com"`
	Password string `json:"password" example:"password123"`
}

// Register mendaftarkan user baru.
// @Summary Mendaftarkan user
// @Description Membuat akun customer baru.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Data registrasi"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/auth/register [post]
func (controller *AuthController) Register(c *gin.Context) {
	var request RegisterRequest
	err := c.ShouldBindJSON(&request)
	email := strings.ToLower(strings.TrimSpace(request.Email))
	_, emailErr := mail.ParseAddress(email)
	if err != nil || emailErr != nil || !strings.Contains(email, "@") || strings.TrimSpace(request.Name) == "" || len(request.Password) < 8 {
		utils.Error(c.Writer, http.StatusBadRequest, "Name, valid email, and password of at least 8 characters are required")
		return
	}

	var existing models.User
	if result := controller.DB.Where("email = ?", email).First(&existing); result.Error == nil {
		utils.Error(c.Writer, http.StatusConflict, "Email is already registered")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to secure password")
		return
	}
	user := models.User{Name: strings.TrimSpace(request.Name), Email: email, PasswordHash: string(passwordHash), Role: "customer"}
	if result := controller.DB.Create(&user); result.Error != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.Success(c.Writer, http.StatusCreated, "Registration successful", userResponse(user))
}

// Login melakukan autentikasi user.
// @Summary Login user
// @Description Menghasilkan access token untuk request terproteksi.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Data login"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/auth/login [post]
func (controller *AuthController) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	result := controller.DB.Where("email = ?", strings.ToLower(strings.TrimSpace(request.Email))).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.Error != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		utils.Error(c.Writer, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.CreateToken(user.ID, user.Role)
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to create access token")
		return
	}
	utils.Success(c.Writer, http.StatusOK, "Login successful", gin.H{"token": token, "user": userResponse(user)})
}

func userResponse(user models.User) gin.H {
	return gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role}
}
