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
	return c.RedirectToRoute("Devices", fiber.Map{}, fiber.StatusMovedPermanently)
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

func DeviceLinePlots(c *fiber.Ctx) error {
	plotType, err := FormatDataType(c.Params("plot_type"))
	if err != nil {
		NotFound(c)
	}

	dateRange, err := strconv.ParseInt(c.Query("range"), 10, 0)
	if err != nil {
		dateRange = 1
	}

	datas, timestamps := database.GetDeviceScrapes(c.Params("deveui"), plotType, int(dateRange))

	return RenderComponent(c,
		com.DevicePlot(
			string(plotType),
			datas,
			timestamps,
		),
	)
}

func DeviceHeatmap(c *fiber.Ctx) error {
	plotType, err := FormatDataType(c.Params("plot_type"))
	if err != nil {
		NotFound(c)
	}

	datas, timestamps := database.GetDeviceScrapesAvg(c.Params("deveui"), plotType)
	return RenderComponent(c,
		com.DeviceHeatmap(
			string(plotType),
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

func Trophies(c *fiber.Ctx) error {
	return RenderComponent(c,
		com.Trophies(
			database.GetTrophies(models.CO2),
			database.GetTrophies(models.Temp),
			database.GetTrophies(models.Humidity),
		),
	)
}
