package brood

type User struct {
	Id              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	NormalizedEmail string `json:"normalized_email"`
	Verified        bool   `json:"verified"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type UserGeneratedToken struct {
	Id        string `json:"access_token"`
	TokenType string `json:"token_type"`
}

type UserToken struct {
	Id        string `json:"id"`
	TokenType string `json:"token_type"`
	CreatedAt string `json:"created_at"`
	UserId    string `json:"user_id"`
	UpdatedAt string `json:"updated_at"`
	Note      string `json:"note"`
	Active    bool   `json:"active"`
}
