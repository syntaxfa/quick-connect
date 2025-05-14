package tokenservice

type TokenGenerateResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
}

type KeyPair struct {
	PrivateKeyHex string
	PublicKeyHex  string
}

type TokenVerifyRequest struct {
	Token string `json:"token"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
