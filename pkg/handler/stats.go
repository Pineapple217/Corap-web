package handler

import (
	"github.com/Pineapple217/Corap-web/pkg/view"

	"github.com/labstack/echo/v4"
)

func (h *Handler) StatsHome(c echo.Context) error {
	stats, err := h.DB.AllStats(c.Request().Context())
	if err != nil {
		return err
	}
	return render(c, view.StatsHome(stats))
}
