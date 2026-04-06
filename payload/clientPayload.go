package payload

// AppType:  WEB | MOBILE | SPA | NATIVE | M2M
// AuthType: EMAIL_PASSWORD | OTP_EMAIL | OTP_PHONE | OAUTH | MAGIC_LINK | SSO

type CreateClientRequest struct {
	TenantID     string   `json:"tenantId"      binding:"required"`
	Name         string   `json:"name"          binding:"required"`
	AppType      string   `json:"appType"       binding:"required"`
	AuthType     string   `json:"authType"      binding:"required"`
	Description  string   `json:"description"`
	RedirectURIs []string `json:"redirectUris"`
}
