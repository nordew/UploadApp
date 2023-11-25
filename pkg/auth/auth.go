package auth

type Authenticator interface {
	// GenerateToken provides opportunity to encrypt access token.
	GenerateToken(options *GenerateTokenClaimsOptions) (string, error)
	// ParseToken provides opportunity to decrypt access token.
	ParseToken(accessToken string) (*ParseTokenClaimsOutput, error)
}

type GenerateTokenClaimsOptions struct {
	UserId   string
	UserName string
}

type ParseTokenClaimsOutput struct {
	UserId   string
	Username string
}
