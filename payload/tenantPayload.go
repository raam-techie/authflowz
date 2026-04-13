package payload

type TenantAuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTenantTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TenantLoginRequest struct {
	AccountID string `json:"accountId" binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

type CreateTenantRequest struct {
	Name       string `json:"name"       binding:"required"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	WebsiteURL string `json:"websiteUrl"`
}
