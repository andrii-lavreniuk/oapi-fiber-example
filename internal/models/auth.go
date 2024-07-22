package models

type Auth struct {
	ID     string
	APIKey string `bun:"api-key"`
}
