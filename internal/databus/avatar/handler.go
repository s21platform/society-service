package avatar

import (
	"context"
	"github.com/s21platform/society-service/internal/config"
)

type Handler struct {
	dbRepo DBRepo
}

func New(dbRepo config.Postgres) *Handler { return &Handler{dbRepo: dbRepo} }

func (h *Handler) Handler(ctx context.Context, in []byte) {}
