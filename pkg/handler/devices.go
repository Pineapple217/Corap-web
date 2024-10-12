package handler

import (
	"net/http"
	"strconv"

	"github.com/Pineapple217/Corap-web/pkg/database"
	"github.com/Pineapple217/Corap-web/pkg/view"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) DevicesHome(c echo.Context) error {
	return render(c, view.DevicesHome())
}

func (h *Handler) DevicesTable(c echo.Context) error {
	devs, err := h.DB.AllDevices(c.Request().Context())
	if err != nil {
		return err
	}
	_ = devs
	return render(c, view.DevicesTable(devs))
}

func (h *Handler) DevicesAnalysisTable(c echo.Context) error {
	devs, err := h.DB.AllDevicesAnalysis(c.Request().Context())
	if err != nil {
		return err
	}
	_ = devs
	return render(c, view.DevicesAnalysisTable(devs))
}

func (h *Handler) DeviceHome(c echo.Context) error {
	dev, err := h.DB.DeviceById(c.Request().Context(), c.Param("deveui"))
	if err == pgx.ErrNoRows {
		return echo.NotFoundHandler(c)
	}
	return render(c, view.DeviceHome(dev))
}

// type historyResponse struct {
// 	Times []time.Time `json:"times"`
// 	Data  []any       `json:"data"`
// }

func (h *Handler) DeviceHistory(c echo.Context) error {
	t, err := database.FormatDataType(c.QueryParam("type"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	dayRange, err := strconv.ParseInt(c.QueryParam("range"), 10, 0)
	if err != nil {
		dayRange = 1
	}
	if dayRange < 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	if dayRange > 90 {
		dayRange = 90
	}

	datas, err := h.DB.DeviceHistory(
		c.Request().Context(),
		c.Param("deveui"),
		t,
		int(dayRange),
	)
	if err == pgx.ErrNoRows {
		return echo.NotFoundHandler(c)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, datas)
}
