package payload

type TenantAuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TenantLoginRequest struct {
	AccountID string `json:"accountId" binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

type CreateTenantRequest struct {
	Name       string  `json:"name"       binding:"required"`
	Email      string  `json:"email"      binding:"required,email"`
	Password   string  `json:"password"   binding:"required,min=8"`
	Phone      string  `json:"phone"`
	Address    Address `json:"address"`
	WebsiteURL string  `json:"websiteUrl"`
}

type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
	Country      string `json:"country"`
}

type Tenant struct {
	ID         string `json:"id"`
	AccountID  string `json:"accountId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Address    string `json:"address,omitempty"`
	WebsiteURL string `json:"websiteUrl,omitempty"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
}
