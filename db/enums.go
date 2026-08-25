package db

// AppType represents the type of application an app client is built for.
type AppType string

const (
	AppTypeWeb    AppType = "WEB"
	AppTypeMobile AppType = "MOBILE"
)

var validAppTypes = map[AppType]struct{}{
	AppTypeWeb:    {},
	AppTypeMobile: {},
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
)

var validSignInMethods = map[SignInMethod]struct{}{
	SignInEmailPassword: {},
	SignInPhonePassword: {},
}

func IsValidSignInMethod(s string) bool {
	_, ok := validSignInMethods[SignInMethod(s)]
	return ok
}
