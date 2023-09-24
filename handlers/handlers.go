package handlers

import (
	"Corap-web/database"
	"Corap-web/models"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

func Home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func Devices(c *fiber.Ctx) error {
	return c.Render("devices", fiber.Map{
		"Devices": database.GetDevices(),
	}, "layouts/main")
}

func Device(c *fiber.Ctx) error {
	device, err := database.GetDevice(c.Params("deveui"))
	if err != nil {
		return c.Status(404).SendFile("./static/private/404.html")
	}
	return c.Render("device", fiber.Map{
		"Device": device,
	}, "layouts/main")
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
		return c.Status(404).SendFile("./static/private/404.html")
	}

	dateRange, err := strconv.ParseInt(c.Query("range"), 10, 0)
	if err != nil {
		dateRange = 1
	}

	datas, timestamps := database.GetDeviceScrapes(c.Params("deveui"), plotType, int(dateRange))
	return c.Render("partials/device_plot", fiber.Map{
		"PlotType": plotTypeStr,
		"PlotData": datas,
		"PlotTime": timestamps,
	})
}

func Jobs(c *fiber.Ctx) error {
	return c.Render("jobs", fiber.Map{
		"Jobs": database.GetSchedulerJobs(),
	}, "layouts/main")
}

func Scrape(c *fiber.Ctx) error {
	return c.Render("scrape", fiber.Map{
		"ScrapeCount":  database.GetScrapeCount(),
		"DatabaseSize": database.GetDatabaseSize(),
		"BatchCount":   database.GetBatchCount(),
		"LastScraped":  database.GetTimeLastScrape().Format("2006-01-02 15:04:05"),
	}, "layouts/main")
}
