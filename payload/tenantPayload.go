package payload

type CreateTenantRequest struct {
	Name       string `json:"name"       binding:"required"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	WebsiteURL string `json:"websiteUrl"`
}
