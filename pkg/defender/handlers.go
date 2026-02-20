package defender

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/labstack/echo/v4"
)

//go:embed templates/defender.gohtml
var defenderTemplateFS embed.FS

var defenderTmpl = template.Must(template.ParseFS(defenderTemplateFS, "templates/defender.gohtml"))

// StatusResponse represents the defender status returned by the API.
type StatusResponse struct {
	Running bool   `json:"running"`
	Config  Config `json:"config"`
}

// speedOption is used by the template to render the speed dropdown.
type speedOption struct {
	Label string
	Value ogame.Speed
}

// pageData holds data passed to the gohtml template.
type pageData struct {
	Running bool
	Config  Config
	Speeds  []speedOption
}

// PageHandler renders the defender configuration HTML page.
func PageHandler(d *Defender) echo.HandlerFunc {
	speeds := []speedOption{
		{"10%", ogame.TenPercent},
		{"20%", ogame.TwentyPercent},
		{"30%", ogame.ThirtyPercent},
		{"40%", ogame.FourtyPercent},
		{"50%", ogame.FiftyPercent},
		{"60%", ogame.SixtyPercent},
		{"70%", ogame.SeventyPercent},
		{"80%", ogame.EightyPercent},
		{"90%", ogame.NinetyPercent},
		{"100%", ogame.HundredPercent},
	}
	return func(c echo.Context) error {
		data := pageData{
			Running: d.IsRunning(),
			Config:  d.GetConfig(),
			Speeds:  speeds,
		}
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		return defenderTmpl.Execute(c.Response().Writer, data)
	}
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
