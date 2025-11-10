package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-se/db"
	"github.com/hewo233/hdu-se/models"
	"github.com/hewo233/hdu-se/shared/consts"
	password "github.com/hewo233/hdu-se/utils"
	"net/http"
)

func CheckUserExistByEmail(email string, c *gin.Context) bool {
	existingUser := models.UserNew()
	result := db.DB.Table(consts.UserTable).Where("email = ?", email).Limit(1).Find(existingUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50001,
			Result: "Database error",
		})
		c.Abort()
		return false
	}
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40002,
			Result: "User with this email already exists",
		})
		c.Abort()
		return false
	}

	return true
}

type registerUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (t registerUserRequest) check() bool {
	if t.Username == "" || t.Email == "" || len(t.Password) < 6 {
		return false
	}

	return true
}

type registerUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterUser Register
func RegisterUser(c *gin.Context) {
	req := registerUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40000,
			Result: "Invalid request data",
		})
		return
	}

	// validate request data
	if !req.check() {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40001,
			Result: "Invalid request data",
		})
		return
	}

	HashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50000,
			Result: "Failed to hash password",
		})
		return
	}

	user := models.UserNew()

	user = &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: HashedPassword,
	}

	// find if user exists
	if !CheckUserExistByEmail(req.Email, c) {
		return
	}

	result := db.DB.Table(consts.UserTable).Create(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50002,
			Result: "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, models.Report{
		Code: http.StatusOK,
		Result: registerUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})

}
