package kis

import (
	"context"
	"kis-flow/config"
)

type Flow interface {
	Run(ctx context.Context) error
	Link(fConfig *config.KisFuncConfig, fParams config.FParam) error
}
