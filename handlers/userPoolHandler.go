package handlers

import (
	"new-auth-service/errutil"
	"new-auth-service/payload"
	"new-auth-service/services"

	"github.com/gin-gonic/gin"
)

func CreateUserPool(c *gin.Context) {
	var req payload.CreateUserPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	pool, err := services.CreateUserPool(c.Request.Context(), &req)
	if err != nil {
		msg := err.Error()
		if msg == "user pool not found" {
			panic(errutil.NotFound(msg))
		}
		if isValidationError(msg) {
			panic(errutil.BadRequest(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "User pool created successfully",
		Data: []any{map[string]any{
			"id":            pool.ID,
			"tenantId":      pool.TenantID,
			"poolId":        pool.PoolID,
			"poolName":      pool.PoolName,
			"description":   pool.Description,
			"signInMethods": pool.SignInMethods,
			"appType":       pool.AppType,
			"isActive":      pool.IsActive,
			"createdAt":     pool.CreatedAt,
		}},
	})
}

func GetUserPool(c *gin.Context) {
	id := c.Query("id")
	tenantID := c.Query("tenantId")

	if id == "" && tenantID == "" {
		panic(errutil.BadRequest("either 'id' or 'tenantId' query param is required"))
	}

	if id != "" {
		pool, err := services.GetUserPoolByID(c.Request.Context(), id)
		if err != nil {
			if err.Error() == "user pool not found" {
				panic(errutil.NotFound(err.Error()))
			}
			panic(errutil.Internal(err.Error()))
		}

		c.JSON(200, payload.SuccessResponse{
			StatusCode: 200,
			Message:    "User pool fetched successfully",
			Data: []any{map[string]any{
				"id":            pool.ID,
				"tenantId":      pool.TenantID,
				"poolId":        pool.PoolID,
				"poolName":      pool.PoolName,
				"description":   pool.Description,
				"signInMethods": pool.SignInMethods,
				"appType":       pool.AppType,
				"isActive":      pool.IsActive,
				"createdAt":     pool.CreatedAt,
			}},
		})
		return
	}

	pools, err := services.GetUserPoolsByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	data := make([]any, 0, len(pools))
	for _, p := range pools {
		data = append(data, map[string]any{
			"id":            p.ID,
			"tenantId":      p.TenantID,
			"poolId":        p.PoolID,
			"poolName":      p.PoolName,
			"description":   p.Description,
			"signInMethods": p.SignInMethods,
			"appType":       p.AppType,
			"isActive":      p.IsActive,
			"createdAt":     p.CreatedAt,
		})
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "User pools fetched successfully",
		Data:       data,
	})
}

func CreateAppClient(c *gin.Context) {
	tenantID := c.Query("tenantId")
	if tenantID == "" {
		panic(errutil.BadRequest("'tenantId' query param is required"))
	}

	var req payload.CreateAppClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(errutil.BadRequest(err.Error()))
	}

	client, clientSecret, err := services.CreateAppClient(c.Request.Context(), tenantID, &req)
	if err != nil {
		msg := err.Error()
		if msg == "user pool not found" {
			panic(errutil.NotFound(msg))
		}
		if isValidationError(msg) {
			panic(errutil.BadRequest(msg))
		}
		panic(errutil.Internal(msg))
	}

	c.JSON(201, payload.SuccessResponse{
		StatusCode: 201,
		Message:    "App client created successfully",
		Data: []any{map[string]any{
			"id":           client.ID,
			"tenantId":     client.TenantID,
			"clientName":   client.ClientName,
			"clientId":     client.ClientID,
			"clientSecret": clientSecret,
			"appType":      client.AppType,
			"isActive":     client.IsActive,
			"createdAt":    client.CreatedAt,
		}},
	})
}

func GetAppClients(c *gin.Context) {
	tenantID := c.Query("tenantId")
	if tenantID == "" {
		panic(errutil.BadRequest("'tenantId' query param is required"))
	}

	clients, err := services.GetAppClientsByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		panic(errutil.Internal(err.Error()))
	}

	data := make([]any, 0, len(clients))
	for _, cl := range clients {
		data = append(data, map[string]any{
			"id":         cl.ID,
			"tenantId":   cl.TenantID,
			"clientName": cl.ClientName,
			"clientId":   cl.ClientID,
			"appType":    cl.AppType,
			"isActive":   cl.IsActive,
			"createdAt":  cl.CreatedAt,
		})
	}

	c.JSON(200, payload.SuccessResponse{
		StatusCode: 200,
		Message:    "App clients fetched successfully",
		Data:       data,
	})
}

// isValidationError returns true for errors that should be surfaced as 400 Bad Request.
func isValidationError(msg string) bool {
	prefixes := []string{
		"invalid app type:",
		"invalid sign-in method:",
		"at least one sign-in method",
	}
	for _, p := range prefixes {
		if len(msg) >= len(p) && msg[:len(p)] == p {
			return true
		}
	}
	return false
}
