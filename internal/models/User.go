package models

type User struct {
	Username       string   `json:"username"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Email_verified bool     `json:"email_verified"`
	Groups         []string `json:"groups"`
	UserID         string   `json:"user_id"`
}

type Claims struct {
	Name            string          `json:"name"`
	Nickname        string          `json:"nickname"`
	Picture         string          `json:"picture"`
	UpdatedAt       string          `json:"updated_at"`
	Email           string          `json:"email"`
	Email_verified  bool            `json:"email_verified"`
	Groups          []string        `json:"groups"`
	Iss             string          `json:"iss"`
	Sub             string          `json:"sub"`
	Aud             string          `json:"aud"`
	Exp             int             `json:"exp"`
	Iat             int             `json:"iat"`
	Nonce           string          `json:"nonce"`
	At_hash         string          `json:"at_hash"`
	FederatedClaims FederatedClaims `json:"federated_claims"`
}

type FederatedClaims struct {
	ConnectorID string `json:"connector_id"`
	UserID      string `json:"user_id"`
}
