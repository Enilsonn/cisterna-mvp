package domain

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RolePipeiro Role = "PIPEIRO"
	RoleCidadao Role = "CIDADAO"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         Role
	IsActive     bool
}

type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  int64
	RefreshExpiry int64
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}

type RefreshToken struct {
	ID       string
	UserID   string
	UserRole Role
	Token    string
	Revoked  bool
}
