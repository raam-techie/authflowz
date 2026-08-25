package payload

type TenantAuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TenantLoginRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	AccountID string `json:"accountId"`
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
