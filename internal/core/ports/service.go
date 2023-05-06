package ports

import "context"

type Service interface {
	GenerateNewsContent(ctx context.Context) error
}
