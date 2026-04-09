package payload

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp"   binding:"required"`
}

type UpdateAttributesRequest struct {
	Attributes map[string]string `json:"attributes"`
}

type CreateUserRequest struct {
	Name       string `json:"name"     binding:"required"`
	Email      string `json:"email"    binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	Role       string `json:"role"`
}
