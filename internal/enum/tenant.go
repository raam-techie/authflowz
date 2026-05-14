package enum

type TenantStatus string

// TenantStatus constants representing the possible states of a tenant account.
const (
	TenantStatusActive  TenantStatus = "ACTIVE"
	TenantStatusHold    TenantStatus = "HOLD"
	TenantStatusBlocked TenantStatus = "BLOCKED"
)

// IsValid checks if the TenantStatus is one of the defined constants.
func (ts TenantStatus) IsValid() bool {
	switch ts {
	case TenantStatusActive, TenantStatusHold, TenantStatusBlocked:
		return true
	default:
		return false
	}
}

// String returns the string representation of the TenantStatus.
func (ts TenantStatus) String() string {
	return string(ts)
}
