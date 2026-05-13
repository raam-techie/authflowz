package handlers

import (
	"net/http"
	"new-auth-service/db"
	"new-auth-service/errutil"
	"new-auth-service/middleware"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	id := c.Query("id")
	tenantID := c.Query("tenantId")

	if id == "" && tenantID == "" {
		panic(errutil.BadRequest("either 'id' or 'tenantId' query param is required"))
	}

	if id != "" {
		user, err := services.GetUserByID(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "user not found" {
				panic(errutil.NotFound(err.Error()))
			}
			panic(errutil.Internal(err.Error()))
		}

		c.JSON(http.StatusOK, payload.SuccessResponse{
			StatusCode: http.StatusOK,
			Message:    "User fetched successfully",
			Data: []any{map[string]any{
				"id":          user.ID,
				"tenantId":    user.TenantID,
				"appClientId": user.AppClientID,
				"userId":      user.UID,
				"name":        user.Name,
				"email":       user.Email,
				"phone":       user.Phone,
				"role":        user.Role,
				"isActive":    user.IsActive,
				"createdAt":   user.CreatedAt,
			}},
		})
		return
	}

	users, err := services.GetUsersByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	data := make([]any, 0, len(users))
	for _, u := range users {
		data = append(data, map[string]any{
			"id":          u.ID,
			"tenantId":    u.TenantID,
			"appClientId": u.AppClientID,
			"userId":      u.UID,
			"name":        u.Name,
			"email":       u.Email,
			"phone":       u.Phone,
			"role":        u.Role,
			"isActive":    u.IsActive,
			"createdAt":   u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Users fetched successfully",
		Data:       data,
	})
}

func appClientFromCtx(c *gin.Context) *db.AppClient {
	a, exists := c.Get(middleware.PoolClientContextKey)
	if !exists {
		panic(errutil.Internal("app client context missing"))
	}
	appClient, ok := a.(*db.AppClient)
	if !ok {
		panic(errutil.Internal("invalid app client context type"))
	}
	return appClient
}

func CreateUser(c *gin.Context) {
	appClient := appClientFromCtx(c)

	var req payload.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	user, err := services.CreateUser(c.Request.Context(), &req, appClient)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	auditUserID := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	} else {
		auditUserID = user.ID
	}
	services.LogAction(services.AuditEntry{
		TenantID:      appClient.TenantID,
		AppClientID:   appClient.ID,
		UserID:        auditUserID,
		Action:        db.AuditActionUserCreate,
		ActorType:     db.ActorTypeUser,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"email": req.Email, "role": req.Role},
	})

	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(http.StatusCreated, payload.SuccessResponse{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data: []any{map[string]any{
			"id":          user.ID,
			"tenantId":    user.TenantID,
			"appClientId": user.AppClientID,
			"userId":      user.UID,
			"name":        user.Name,
			"email":       user.Email,
			"phone":       user.Phone,
			"role":        user.Role,
			"isActive":    user.IsActive,
			"createdAt":   user.CreatedAt,
		}},
	})
}

func SendEmailOTP(c *gin.Context) {
	appClient := appClientFromCtx(c)

	var req payload.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	err := services.SendEmailOTP(c.Request.Context(), &req, appClient)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		TenantID:      appClient.TenantID,
		AppClientID:   appClient.ID,
		Action:        db.AuditActionOTPSend,
		ActorType:     db.ActorTypeUser,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"email": req.Email},
	})

	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "OTP sent successfully",
		Data:       []any{},
	})
}

func VerifyEmailOTP(c *gin.Context) {
	appClient := appClientFromCtx(c)

	var req payload.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	token, err := services.VerifyEmailOTP(c.Request.Context(), &req, appClient)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		TenantID:      appClient.TenantID,
		AppClientID:   appClient.ID,
		Action:        db.AuditActionOTPVerify,
		ActorType:     db.ActorTypeUser,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"email": req.Email},
	})

	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       []any{token},
	})
}

func UserLogin(c *gin.Context) {
	appClient := appClientFromCtx(c)

	var req payload.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	tokens, err := services.LoginUser(c.Request.Context(), &req, appClient)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		TenantID:      appClient.TenantID,
		AppClientID:   appClient.ID,
		Action:        db.AuditActionUserLogin,
		ActorType:     db.ActorTypeUser,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"email": req.Email},
	})

	if err != nil {
		msg := err.Error()
		if msg == "invalid credentials" || msg == "user account is not active" || msg == "email/password auth is not enabled for this pool" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       []any{tokens},
	})
}

func UpdateUserAttributes(c *gin.Context) {
	appClient := appClientFromCtx(c)
	userID := c.Param("userId")

	var req payload.UpdateAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	if len(req.Attributes) == 0 {
		panic(errutil.BadRequest("attributes must not be empty"))
	}

	err := services.UpdateUserAttributes(c.Request.Context(), userID, &req, appClient)

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		TenantID:      appClient.TenantID,
		AppClientID:   appClient.ID,
		UserID:        userID,
		Action:        db.AuditActionUserAttributeUpdate,
		ActorType:     db.ActorTypeUser,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
		Metadata:      map[string]any{"attributeCount": len(req.Attributes)},
	})

	if err != nil {
		if err.Error() == "user not found" {
			panic(errutil.NotFound(err.Error()))
		}
		panic(errutil.Internal(err.Error()))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
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

	auditStatus := db.AuditStatusSuccess
	auditReason := ""
	if err != nil {
		auditStatus = db.AuditStatusFailure
		auditReason = err.Error()
	}
	services.LogAction(services.AuditEntry{
		Action:        db.AuditActionTokenRefresh,
		ActorType:     db.ActorTypeSystem,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		Status:        auditStatus,
		FailureReason: auditReason,
	})

	if err != nil {
		msg := err.Error()
		if msg == "invalid refresh token" || msg == "refresh token expired" ||
			msg == "user account is not active" || msg == "tenant account is not active" {
			panic(errutil.Unauthorized(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(http.StatusOK, payload.SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Token refreshed successfully",
		Data:       []any{tokens},
	})
}
