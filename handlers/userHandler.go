package handlers

import (
	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req payload.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	user, err := services.CreateUser(c.Request.Context(), &req)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "User created successfully",
		Data: []any{map[string]any{
			"id":         user.ID,
			"tenantId":   user.TenantID,
			"projectId":  user.ProjectID,
			"userId":     user.UID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"department": user.Department,
			"role":       user.Role,
			"isActive":   user.IsActive,
			"createdAt":  user.CreatedAt,
		}},
	})
}

func SendEmailOTP(c *gin.Context) {
	var req payload.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	if err := services.SendEmailOTP(c.Request.Context(), &req); err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "OTP sent successfully",
		Data:       []any{},
	})
}

func VerifyEmailOTP(c *gin.Context) {
	var req payload.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.VerifyEmailOTP(c.Request.Context(), &req)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{token},
	})
}

func GoogleOAuthLogin(c *gin.Context) {
	var req payload.GoogleOAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.GoogleOAuthLogin(c.Request.Context(), &req)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{token},
	})
}

func UserLogin(c *gin.Context) {
	var req payload.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.LoginUser(c.Request.Context(), &req)
	if err != nil {
		msg := err.Error()
		if msg == "invalid credentials" || msg == "user account is not active" || msg == "project is not active" || msg == "email/password auth is not enabled for this project" {
			panic(errutil.Unauthorized(msg))
		}
		if msg == "project not found" {
			panic(errutil.NotFound(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{token},
	})
}
