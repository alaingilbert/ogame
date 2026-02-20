package defender

import (
	"net/http"
	"strconv"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/labstack/echo/v4"
)

// StatusResponse represents the defender status returned by the API.
type StatusResponse struct {
	Running bool   `json:"running"`
	Config  Config `json:"config"`
}

// GetStatusHandler returns the current defender status and config.
func GetStatusHandler(d *Defender) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, StatusResponse{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
		})
	}
}

// ConfigRequest represents the JSON body for configuring the defender.
type ConfigRequest struct {
	Enabled                     *bool    `json:"enabled"`
	CheckIntervalSec            *int     `json:"checkIntervalSec"`
	EvacSpeed                   *float64 `json:"evacSpeed"`
	EvacSecondsBefore           *int     `json:"evacSecondsBefore"`
	RecallSecondsAfter          *int     `json:"recallSecondsAfter"`
	DumpResourcesToDefense      *bool    `json:"dumpResourcesToDefense"`
	DumpResourcesToDefenseMoons *bool    `json:"dumpResourcesToDefenseMoons"`
	OriginRandom                *bool    `json:"originRandom"`
	OriginCelestialID           *int64   `json:"originCelestialID"`
}

// SetConfigHandler updates the defender configuration.
func SetConfigHandler(d *Defender) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req ConfigRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		cfg := d.GetConfig()

		if req.Enabled != nil {
			cfg.Enabled = *req.Enabled
		}
		if req.CheckIntervalSec != nil {
			cfg.CheckIntervalSec = *req.CheckIntervalSec
		}
		if req.EvacSpeed != nil {
			cfg.EvacSpeed = ogame.Speed(*req.EvacSpeed)
		}
		if req.EvacSecondsBefore != nil {
			cfg.EvacSecondsBefore = *req.EvacSecondsBefore
		}
		if req.RecallSecondsAfter != nil {
			cfg.RecallSecondsAfter = *req.RecallSecondsAfter
		}
		if req.DumpResourcesToDefense != nil {
			cfg.DumpResourcesToDefense = *req.DumpResourcesToDefense
		}
		if req.DumpResourcesToDefenseMoons != nil {
			cfg.DumpResourcesToDefenseMoons = *req.DumpResourcesToDefenseMoons
		}
		if req.OriginRandom != nil {
			cfg.Origin.Random = *req.OriginRandom
		}
		if req.OriginCelestialID != nil {
			cfg.Origin.CelestialID = ogame.CelestialID(*req.OriginCelestialID)
		}

		d.SetConfig(cfg)

		return c.JSON(http.StatusOK, StatusResponse{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
		})
	}
}

// StartHandler starts the defender.
func StartHandler(d *Defender) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg := d.GetConfig()
		cfg.Enabled = true
		d.SetConfig(cfg)
		d.Start()
		return c.JSON(http.StatusOK, StatusResponse{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
		})
	}
}

// StopHandler stops the defender.
func StopHandler(d *Defender) echo.HandlerFunc {
	return func(c echo.Context) error {
		d.Stop()
		return c.JSON(http.StatusOK, StatusResponse{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
		})
	}
}

// SetOriginHandler sets the origin for attack checks.
func SetOriginHandler(d *Defender) echo.HandlerFunc {
	return func(c echo.Context) error {
		mode := c.QueryParam("mode")
		cfg := d.GetConfig()

		switch mode {
		case "random":
			cfg.Origin.Random = true
		default:
			celestialIDStr := c.QueryParam("celestialID")
			if celestialIDStr == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "celestialID required when mode is not random"})
			}
			id, err := strconv.ParseInt(celestialIDStr, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid celestialID"})
			}
			cfg.Origin.Random = false
			cfg.Origin.CelestialID = ogame.CelestialID(id)
		}

		d.SetConfig(cfg)
		return c.JSON(http.StatusOK, StatusResponse{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
		})
	}
}
