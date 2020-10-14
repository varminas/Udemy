package model

// https://charnnarong.github.io/json-to-go-struct/
// AccessTokenResponse response from Authorization server after access token exchange
type AccessTokenResponse struct {
	AccessToken	string	`json:"access_token"`
	ExpiresIn	int	`json:"expires_in"`
	RefreshExpiresIn	int	`json:"refresh_expires_in"`
	RefreshToken	string	`json:"refresh_token"`
	TokenType	string	`json:"token_type"`
	NotBeforePolicy	int	`json:"not-before-policy"`
	SessionState	string	`json:"session_state"`
	Scope	string	`json:"scope"`
  }
  