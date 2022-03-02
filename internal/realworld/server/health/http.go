package health

import (
	"context"
	"github.com/shencw/kratos-realworld/pkg/transport/http/health"
	"net/http"
)

func NewHttpHealthHandler() http.Handler {
	healthHandler := health.NewHandler()
	healthHandler.AddChecker("ping", func(ctx context.Context) error {
		return nil
	})
	return healthHandler
}
