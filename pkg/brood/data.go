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
