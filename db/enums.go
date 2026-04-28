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

var validAppTypes = map[AppType]struct{}{
	AppTypeWeb:    {},
	AppTypeMobile: {},
	AppTypeSPA:    {},
	AppTypeNative: {},
	AppTypeM2M:    {},
}

func IsValidAppType(s string) bool {
	_, ok := validAppTypes[AppType(s)]
	return ok
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

var validSignInMethods = map[SignInMethod]struct{}{
	SignInEmailPassword: {},
	SignInPhonePassword: {},
	SignInOTPEmail:      {},
	SignInOTPPhone:      {},
	SignInMagicLink:     {},
	SignInSSO:           {},
}

func IsValidSignInMethod(s string) bool {
	_, ok := validSignInMethods[SignInMethod(s)]
	return ok
}
