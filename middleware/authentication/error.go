package authentication

import "errors"

const (
	invalidToken       = "invalid token"
	invalidPathParam   = "invalid path param userId"
	unauthorizedAccess = "unauthorized access"
)

var (
	noUserTokenError         = errors.New("failed to get user token")
	invalidTokenClaimsError  = errors.New("invalid token claims")
	invalidSubjectClaimError = errors.New("invalid subject claim")
)
