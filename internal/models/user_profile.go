package models

type UserProfile struct {
	UserID    int64  `bun:"user_id,pk"`
	FirstName string `bun:"first_name"`
	LastName  string `bun:"last_name"`
	Phone     string `bun:"phone"`
	Address   string `bun:"address"`
	City      string `bun:"city"`
}
