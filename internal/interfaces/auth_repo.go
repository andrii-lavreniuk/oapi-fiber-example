package interfaces

import (
	"context"
)

type AuthRepo interface {
	Exists(context.Context, string) (bool, error)
}
