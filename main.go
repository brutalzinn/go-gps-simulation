package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/brutalzinn/go-gps-simulation/adbhelper"
	"github.com/brutalzinn/go-gps-simulation/config"
	"github.com/brutalzinn/go-gps-simulation/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var adbHelper *adbhelper.AdbHelper

func main() {
	app := fiber.New()
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Error loading config: %v", err)
	}
	adbHelper = adbhelper.New(config.AdbPath)
	app.Static("/", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"GoogleMapsAPIKey": config.GoogleMapsAPIKey,
			"DefaultDevice":    config.DefaultDevice,
		})
	})
	app.Post("/update", func(c *fiber.Ctx) error {
		var payload models.RequestPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		go simulateGPS(payload.PointA, payload.PointB, payload.Device)
		return c.SendString("Simulation started")
	})
	app.Get("/devices", func(c *fiber.Ctx) error {
		output, err := adbHelper.SendCommandWithResults("devices")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error listing devices")
		}
		lines := strings.Split(string(output), "\n")
		var devices []string
		for _, line := range lines {
			if strings.Contains(line, "\tdevice") {
				devices = append(devices, strings.Fields(line)[0])
			}
		}

		return c.JSON(devices)
	})
	if err := adbHelper.SetupADB(config.AdbBaseURL); err != nil {
		logrus.Fatalf("Error setting up ADB: %v", err)
	}
	logrus.Info("Server started at :%d", config.Port)
	logrus.Info(app.Listen(fmt.Sprintf(":%d", config.Port)))
}

func simulateGPS(pointA, pointB models.Coordinates, device string) {
	steps := 10 /// gamb: number of steps
	latStep := (pointB.Lat - pointA.Lat) / float64(steps)
	lngStep := (pointB.Lng - pointA.Lng) / float64(steps)

	for i := 0; i <= steps; i++ {
		lat := pointA.Lat + latStep*float64(i)
		lng := pointA.Lng + lngStep*float64(i)
		adbHelper.SendCommand("-s " + device + " emu geo fix " + fmt.Sprintf("%f", lng) + " " + fmt.Sprintf("%f", lat))
		time.Sleep(1 * time.Second) /// delay between changes
	}
}
