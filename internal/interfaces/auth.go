package interfaces

import (
	"context"
)

type Auth interface {
	ValidateAPIKey(ctx context.Context, apikey string) (bool, error)
}
