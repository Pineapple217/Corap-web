package handlers

import (
	com "Corap-web/components"
	"Corap-web/database"
	"Corap-web/models"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendFile("./static/private/404.html")
}

func Home(c *fiber.Ctx) error {
	return c.Redirect("devices", fiber.StatusMovedPermanently)
	// return RenderComponent(c, com.Main())
}

func Devices(c *fiber.Ctx) error {
	return RenderComponent(c,
		com.Devices(),
	)
}

func DevicesTable(c *fiber.Ctx) error {
	return RenderComponent(c,
		com.DeviceTable(
			database.GetDevices(),
		),
	)
}

func Device(c *fiber.Ctx) error {
	device, err := database.GetDevice(c.Params("deveui"))
	if err != nil {
		return NotFound(c)
	}
	return RenderComponent(c,
		com.Device(device),
	)
}

func DevicePlots(c *fiber.Ctx) error {
	var plotType models.PlotType
	var plotTypeStr string
	switch c.Params("plot_type") {
	case string(models.Temp):
		plotType = models.Temp
		plotTypeStr = "Temperature"
	case string(models.CO2):
		plotType = models.CO2
		plotTypeStr = "Co2"
	case string(models.Humidity):
		plotType = models.Humidity
		plotTypeStr = "Humidity"
	default:
		return NotFound(c)
	}

	dateRange, err := strconv.ParseInt(c.Query("range"), 10, 0)
	if err != nil {
		dateRange = 1
	}

	datas, timestamps := database.GetDeviceScrapes(c.Params("deveui"), plotType, int(dateRange))
	return RenderComponent(c,
		com.DevicePlot(
			plotTypeStr,
			datas,
			timestamps,
		),
	)
}

func Jobs(c *fiber.Ctx) error {
	return RenderComponent(c,
		com.Jobs(
			database.GetSchedulerJobs(),
		),
	)
}

func Scrape(c *fiber.Ctx) error {
	return RenderComponent(c,
		com.Scrape(
			database.GetTimeLastScrape(),
			database.GetDatabaseSize(),
			database.GetScrapeCount(),
			database.GetBatchCount(),
		),
	)
}
