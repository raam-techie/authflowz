package db

// AppType represents the type of application an app client is built for.
type AppType string

const (
	AppTypeWeb    AppType = "WEB"
	AppTypeMobile AppType = "MOBILE"
	AppTypeSPA    AppType = "SPA"
	AppTypeNative AppType = "NATIVE"
	AppTypeM2M    AppType = "M2M"
)

// IsValid returns true if the AppType is one of the defined constants.
func (a AppType) IsValid() bool {
	switch a {
	case AppTypeWeb, AppTypeMobile, AppTypeSPA, AppTypeNative, AppTypeM2M:
		return true
	default:
		return false
	}
}

// SignInMethod represents an authentication method enabled on a user pool.
type SignInMethod string

const (
	SignInEmailPassword SignInMethod = "EMAIL_PASSWORD"
	SignInPhonePassword SignInMethod = "PHONE_PASSWORD"
	SignInOTPEmail      SignInMethod = "OTP_EMAIL"
	SignInOTPPhone      SignInMethod = "OTP_PHONE"
	SignInMagicLink     SignInMethod = "MAGIC_LINK"
	SignInSSO           SignInMethod = "SSO"
)

// IsValid returns true if the SignInMethod is one of the defined constants.
func (s SignInMethod) IsValid() bool {
	switch s {
	case SignInEmailPassword, SignInPhonePassword, SignInOTPEmail, SignInOTPPhone, SignInMagicLink, SignInSSO:
		return true
	default:
		return false
	}
}
