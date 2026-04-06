package payload

type UserLoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SendOTPRequest struct {
	ProjectID string `json:"projectId" binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
}

type VerifyOTPRequest struct {
	ProjectID string `json:"projectId" binding:"required"`
	Email     string `json:"email"     binding:"required,email"`
	OTP       string `json:"otp"       binding:"required"`
}

type GoogleOAuthLoginRequest struct {
	ProjectID string `json:"projectId" binding:"required"`
	IDToken   string `json:"idToken"   binding:"required"`
}

type CreateUserRequest struct {
	TenantID   string `json:"tenantId"   binding:"required"`
	ProjectID  string `json:"projectId"  binding:"required"`
	Name       string `json:"name"       binding:"required"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	Role       string `json:"role"`
}
