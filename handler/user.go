package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-se/db"
	"github.com/hewo233/hdu-se/models"
	"github.com/hewo233/hdu-se/shared/consts"
	"github.com/hewo233/hdu-se/utils/jwt"
	"github.com/hewo233/hdu-se/utils/password"
	"net/http"
	"strconv"
)

// CheckUserExistByEmail Check if user exists by email
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

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40004,
			Result: "invalid request data",
		})
		c.Abort()
		return
	}

	user := models.UserNew()

	result := db.DB.Table(consts.UserTable).Where("email = ?", req.Email).First(user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusBadRequest, models.Report{
				Code:   40005,
				Result: "user not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Report{
				Code:   50003,
				Result: "database error: " + result.Error.Error(),
			})
		}
		c.Abort()
		return
	}

	if err := password.CheckHashed(req.Password, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, models.Report{
			Code:   40006,
			Result: "incorrect password or email",
		})
		c.Abort()
		return
	}

	strID := strconv.Itoa(int(user.ID))

	jwtToken, err := jwt.GenerateJWT(strID, consts.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Report{
			Code:   50005,
			Result: "failed to generate jwt token",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, models.Report{
		Code: 20000,
		Result: UserLoginResponse{
			User:  *user,
			Token: jwtToken,
		},
	})

}

// CheckUserAuth Check user auth
func CheckUserAuth(id uint, c *gin.Context) bool {
	// get jwt id
	jwtID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Report{
			Code:   40100,
			Result: "unauthorized",
		})
		c.Abort()
		return false
	}

	if jwtID != strconv.Itoa(int(id)) {
		c.JSON(http.StatusUnauthorized, models.Report{
			Code:   40101,
			Result: "unauthorized",
		})
		c.Abort()
		return false
	}

	return true
}

// GetUserInfoByID /user/:id
func GetUserInfoByID(c *gin.Context) {
	userID := c.Param("id")

	user := models.UserNew()
	result := db.DB.Table(consts.UserTable).Where("id = ?", userID).First(user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusBadRequest, models.Report{
				Code:   40007,
				Result: "user not found",
			})
		}
		c.Abort()
		return
	}

	if !CheckUserAuth(user.ID, c) {
		return
	}

	c.JSON(http.StatusOK, models.Report{
		Code:   20000,
		Result: user,
	})
}

// GetUserInfoByEmail /user?email=...
func GetUserInfoByEmail(c *gin.Context) {
	email := c.Query("email")

	user := models.UserNew()
	result := db.DB.Table(consts.UserTable).Where("email = ?", email).First(user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusBadRequest, models.Report{
				Code:   40008,
				Result: "user not found",
			})
			c.Abort()
			return
		}
	}

	if !CheckUserAuth(user.ID, c) {
		return
	}

	c.JSON(http.StatusOK, models.Report{
		Code:   20000,
		Result: user,
	})
}
