package gothquickbooks

// Scopes defines the available permission scopes for QuickBooks API and OpenID Connect
const (
	// QuickBooks API scopes
	ScopeAccounting = "com.intuit.quickbooks.accounting"
	ScopePayment    = "com.intuit.quickbooks.payment"

	// OpenID Connect scopes
	ScopeOpenID  = "openid"
	ScopeProfile = "profile"
	ScopeEmail   = "email"
	ScopePhone   = "phone"
	ScopeAddress = "address"
)
