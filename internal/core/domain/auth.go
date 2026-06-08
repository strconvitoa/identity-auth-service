package domain

type Auth struct {
	ID        string `json:"id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	ExpiresAt string `jaon:"expires_at"`
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
}
