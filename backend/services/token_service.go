package services

// TokenService exposes token helpers to the frontend via Wails binding.
type TokenService struct{}

// NewTokenService creates a new TokenService instance.
func NewTokenService() *TokenService {
	return &TokenService{}
}

// HasOAuthToken checks whether an OAuth token exists in the system keyring.
func (s *TokenService) HasOAuthToken() (bool, error) {
	return HasOAuthToken()
}

// BeginOAuthLogin starts OAuth2 Authorization Code login and stores token in keyring.
func (s *TokenService) BeginOAuthLogin(clientID string, scope string) (bool, error) {
	if err := beginOAuthLogin(clientID, scope); err != nil {
		return false, err
	}
	return true, nil
}

// ClearOAuthToken removes OAuth token from keyring.
func (s *TokenService) ClearOAuthToken() error {
	return ClearOAuthToken()
}
