package extauthapi

const (
	// SessionCookieName is a name of HTTP cookie with session token.
	SessionCookieName = "__Secure-authKey"
	// CSRFTokenCookieName is a name of HTTP cookie with CSRF token.
	CSRFTokenCookieName = "__Secure-CSRFToken" //nolint:gosec
	// CSRFTokenHeaderName is a name of HTTP header with CSRF token bound to session token.
	CSRFTokenHeaderName = "X-CSRFTokenBound" //nolint:gosec
)
