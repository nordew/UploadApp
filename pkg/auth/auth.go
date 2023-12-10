package auth

type Authenticator interface {
	// GenerateTokens provides opportunity to encrypt access & refresh token.
	GenerateTokens(options *GenerateTokenClaimsOptions) (string, string, error)

	// GenerateRefreshToken generates refresh token
	GenerateRefreshToken(id string, role string) (string, error)

	// ParseToken provides opportunity to decrypt access token.
	ParseToken(accessToken string) (*ParseTokenClaimsOutput, error)
}

type GenerateTokenClaimsOptions struct {
	UserId string `json:"sub"`
	Role   string `json:"role""`
}

type ParseTokenClaimsOutput struct {
	UserId   string
	Username string
}
