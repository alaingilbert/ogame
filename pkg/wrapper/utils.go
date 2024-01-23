package wrapper

import (
	"errors"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// GetFleetSpeedForMission ...
func GetFleetSpeedForMission(serverData gameforge.ServerData, missionID ogame.MissionID) int64 {
	if missionID == ogame.Attack ||
		missionID == ogame.GroupedAttack ||
		missionID == ogame.Destroy ||
		missionID == ogame.MissileAttack ||
		missionID == ogame.RecycleDebrisField {
		return serverData.SpeedFleetWar
	}
	return serverData.SpeedFleetPeaceful
}

// ConvertIntoCoordinate helper that turns any type into a coordinate
func ConvertIntoCoordinate(w Wrapper, v IntoCoordinate) (ogame.Coordinate, error) {
	switch vv := v.(type) {
	case string:
		return ogame.ParseCoord(vv)
	case ogame.Coordinate:
		return vv, nil
	case ogame.Celestial:
		return vv.GetCoordinate(), nil
	case ogame.Planet:
		return vv.GetCoordinate(), nil
	case ogame.Moon:
		return vv.GetCoordinate(), nil
	case ogame.CelestialID, ogame.PlanetID, ogame.MoonID:
		c := w.GetCachedCelestial(vv)
		if c == nil {
			return ogame.Coordinate{}, errors.New("celestial not found")
		}
		return c.GetCoordinate(), nil
	default:
		return ogame.Coordinate{}, errors.New("invalid type")
	}
}

// ConvertIntoPlanet helper that turns any type into a Planet
func ConvertIntoPlanet(w Wrapper, v IntoPlanet) (Planet, error) {
	if c := w.GetCachedCelestial(v); c != nil {
		if p, ok := c.(Planet); ok {
			return p, nil
		}
		return Planet{}, errors.New("not a planet")
	}
	return Planet{}, errors.New("planet not found")
}

// ConvertIntoMoon helper that turns any type into a Moon
func ConvertIntoMoon(w Wrapper, v IntoMoon) (Moon, error) {
	if c := w.GetCachedCelestial(v); c != nil {
		if m, ok := c.(Moon); ok {
			return m, nil
		}
		return Moon{}, errors.New("not a moon")
	}
	return Moon{}, errors.New("moon not found")
}
