package models

type UserData struct {
	UserID int64  `bun:"user_id,pk"`
	School string `bun:"school"`
}
