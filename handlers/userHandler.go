package handlers

import (
	"new-auth-service/db"
	"new-auth-service/errutil"
	"new-auth-service/middleware"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func projectFromCtx(c *gin.Context) *db.Project {
	p, exists := c.Get(middleware.ProjectContextKey)
	if !exists {
		panic(errutil.Internal("project context missing"))
	}
	project, ok := p.(*db.Project)
	if !ok {
		panic(errutil.Internal("invalid project context type"))
	}
	return project
}

func CreateUser(c *gin.Context) {
	project := projectFromCtx(c)

	var req payload.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	user, err := services.CreateUser(c.Request.Context(), &req, project)
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
	project := projectFromCtx(c)

	var req payload.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	if err := services.SendEmailOTP(c.Request.Context(), &req, project); err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "OTP sent successfully",
		Data:       []any{},
	})
}

func VerifyEmailOTP(c *gin.Context) {
	project := projectFromCtx(c)

	var req payload.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.VerifyEmailOTP(c.Request.Context(), &req, project)
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
	project := projectFromCtx(c)

	var req payload.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	tokens, err := services.LoginUser(c.Request.Context(), &req, project)
	if err != nil {
		msg := err.Error()
		if msg == "invalid credentials" || msg == "user account is not active" || msg == "email/password auth is not enabled for this project" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data:       []any{tokens},
	})
}

func UpdateUserAttributes(c *gin.Context) {
	project := projectFromCtx(c)
	userID := c.Param("userId")

	var req payload.UpdateAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	if len(req.Attributes) == 0 {
		panic(errutil.BadRequest("attributes must not be empty"))
	}

	if err := services.UpdateUserAttributes(c.Request.Context(), userID, &req, project); err != nil {
		if err.Error() == "user not found" {
			panic(errutil.NotFound(err.Error()))
		}
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Attributes updated successfully",
		Data:       []any{},
	})
}

func RefreshToken(c *gin.Context) {
	var req payload.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	tokens, err := services.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		msg := err.Error()
		if msg == "invalid refresh token" || msg == "refresh token expired" ||
			msg == "user account is not active" || msg == "tenant account is not active" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "Token refreshed successfully",
		Data:       []any{tokens},
	})
}
