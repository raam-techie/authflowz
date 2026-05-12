package payload

type CreateUserPoolRequest struct {
	TenantID      string   `json:"tenantId"      binding:"required"`
	PoolName      string   `json:"poolName"      binding:"required"`
	SignInMethods []string `json:"signInMethods" binding:"required"`
	AppType       string   `json:"appType"       binding:"required"`
	Description   string   `json:"description"`
}

type CreateAppClientRequest struct {
	ClientName string `json:"clientName" binding:"required"`
	AppType    string `json:"appType"    binding:"required"`
}
