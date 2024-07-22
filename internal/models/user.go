package models

type User struct {
	ID       int32  `bun:",pk,autoincrement"`
	Username string `bun:"username"`

	Profile *UserProfile `bun:"rel:has-one,join:id=user_id"`
	Data    *UserData    `bun:"rel:has-one,join:id=user_id"`
}
