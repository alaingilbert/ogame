package v10

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
)

type (
	GalaxyInfosJson struct {
		Success bool       `json:"success"`
		Token   string     `json:"newAjaxToken"`
		System  SystemJson `json:"system"`
	}

	SystemJson struct {
		Galaxy            CInt64        `json:"galaxy"`
		System            CInt64        `json:"system"`
		DeuteriumInDebris bool          `json:"deuteriumInDebris"`
		Content           []ContentJson `json:"galaxyContent"`
	}

	ContentJson struct {
		Galaxy   CInt64      `json:"galaxy"`
		System   CInt64      `json:"system"`
		Position CInt64      `json:"position"`
		Planets  PlanetsList `json:"planets"`
		Player   PlayerJson  `json:"player"`
		Filters  string      `json:"positionFilters"`
	}

	PlanetsList []PlanetJson

	PlanetJson struct {
		PlayerID      CInt64        `json:"playerId"`
		PlanetID      CInt64        `json:"planetId"`
		PlanetName    string        `json:"planetName"`
		Image         string        `json:"imageInformation"`
		PlanetType    CInt64        `json:"planetType"`
		IsDestroyed   bool          `json:"isDestroyed"`
		Size          CInt64        `json:"size"`
		RequiredShips CInt64        `json:"requiredShips"`
		Resources     ResourcesJson `json:"resources"`
		Activity      ActivityJson  `json:"activity"`
	}

	ResourcesJson struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}

	ActivityJson struct {
		Minutes int64
	}

	PlayerJson struct {
		PlayerID          CInt64   `json:"playerId"`
		PlayerName        string   `json:"playerName"`
		AllianceID        CInt64   `json:"allianceId"`
		AllianceName      string   `json:"allianceName"`
		AllianceTag       string   `json:"allianceTag"`
		IsAllianceMember  bool     `json:"isAllianceMember"`
		PositionPlayer    CInt64   `json:"highscorePositionPlayer"`
		PositionAlliance  CInt64   `json:"highscorePositionAlliance"`
		IsAdmin           bool     `json:"isAdmin"`
		IsBanned          bool     `json:"isBanned"`
		IsOnVacation      bool     `json:"isOnVacation"`
		IsNewbie          bool     `json:"isNewbie"`
		IsStrong          bool     `json:"isStrong"`
		IsHonorableTarget bool     `json:"isHonorableTarget"`
		IsInactive        bool     `json:"isInactive"`
		Rank              RankJson `json:"rank"`
	}

	RankJson struct {
		HasRank   bool   `json:"hasRank"`
		RankTitle string `json:"rankTitle"`
		RankClass string `json:"rankClass"`
	}

	CInt64 int64
)

func (c *CInt64) UnmarshalJSON(d []byte) error {
	var tmp int64
	if err := json.Unmarshal(d, &tmp); err == nil {
		*c = CInt64(tmp)
		return nil
	}
	var str string
	if err := json.Unmarshal(d, &str); err != nil {
		return err
	}
	tmp = utils.ParseInt(str)
	*c = CInt64(tmp)
	return nil
}

func (c CInt64) Int64() int64 {
	return int64(c)
}

func (p *PlanetsList) UnmarshalJSON(d []byte) error {
	if d[0] == '[' {
		if err := json.Unmarshal(d, (*[]PlanetJson)(p)); err != nil {
			return err
		}
	} else {
		var tmp PlanetJson
		if err := json.Unmarshal(d, &tmp); err != nil {
			return err
		}
		*p = []PlanetJson{tmp}
		return nil
	}
	return nil
}

func (r *ResourcesJson) UnmarshalJSON(d []byte) error {
	var tmp struct {
		Metal struct {
			Amount CInt64 `json:"amount"`
		} `json:"metal"`
		Crystal struct {
			Amount CInt64 `json:"amount"`
		} `json:"crystal"`
		Deuterium struct {
			Amount CInt64 `json:"amount"`
		} `json:"deuterium"`
	}

	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}

	*r = ResourcesJson{
		Metal:     int64(tmp.Metal.Amount),
		Crystal:   int64(tmp.Crystal.Amount),
		Deuterium: int64(tmp.Deuterium.Amount),
	}

	return nil
}

func (a *ActivityJson) UnmarshalJSON(d []byte) error {

	var min int64
	var tmp map[string]any
	var hour float64 = 60

	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}

	if tmp["showActivity"] == false {
		min = 0
	} else {
		if tmp["showMinutes"] == true && tmp["showActivity"] == hour {
			s := fmt.Sprintf("%v", tmp["idleTime"])
			min = utils.ParseInt(s)
		} else {
			min = 15
		}
	}

	*a = ActivityJson{
		Minutes: min,
	}
	return nil
}

func extractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (ogame.SystemInfos, error) {
	var res ogame.SystemInfos
	// Check for login
	var check struct {
		Success bool
	}
	if err := json.Unmarshal(pageHTML, &check); err != nil {
		return res, ogame.ErrNotLogged
	}
	if !check.Success {
		return res, errors.New("error in retrieve galaxy infos")
	}
	// Real extraction
	var tmp GalaxyInfosJson
	if err := json.Unmarshal(pageHTML, &tmp); err != nil {
		return res, err
	}
	res.OverlayToken = tmp.Token
	res.SetGalaxy(tmp.System.Galaxy.Int64())
	res.SetSystem(tmp.System.System.Int64())

	for i, pos := range tmp.System.Content {

		if pos.Position.Int64() == 16 {
			res.ExpeditionDebris.Metal = pos.Planets[0].Resources.Metal
			res.ExpeditionDebris.Crystal = pos.Planets[0].Resources.Crystal
			res.ExpeditionDebris.Deuterium = pos.Planets[0].Resources.Deuterium
			res.ExpeditionDebris.PathfindersNeeded = pos.Planets[0].RequiredShips.Int64()
			continue
		}
		// TODO: manage asteroids and events in P17
		if pos.Position.Int64() == 17 {
			continue
		}
		if len(pos.Planets) == 0 {
			continue
		}
		player := pos.Player
		planetInfos := new(ogame.PlanetInfos)

		for _, planet := range pos.Planets {

			// Planet
			if planet.PlanetType.Int64() == ogame.PlanetType.Int64() {
				// Generic infos
				planetInfos.ID = planet.PlanetID.Int64()
				planetInfos.Name = planet.PlanetName
				planetInfos.Img = planet.Image

				planetInfos.Inactive = player.IsInactive
				planetInfos.StrongPlayer = player.IsStrong
				planetInfos.Newbie = player.IsNewbie
				planetInfos.Vacation = player.IsOnVacation
				planetInfos.HonorableTarget = player.IsHonorableTarget
				planetInfos.Administrator = player.IsAdmin
				planetInfos.Banned = player.IsBanned
				planetInfos.Destroyed = planet.IsDestroyed

				planetInfos.Activity = planet.Activity.Minutes

				planetInfos.Coordinate.Galaxy = pos.Galaxy.Int64()
				planetInfos.Coordinate.System = pos.System.Int64()
				planetInfos.Coordinate.Position = pos.Position.Int64()
				planetInfos.Coordinate.Type = ogame.PlanetType

				planetInfos.Date = time.Now()

				// Player
				planetInfos.Player.ID = player.PlayerID.Int64()
				planetInfos.Player.Name = player.PlayerName
				planetInfos.Player.Rank = player.PositionPlayer.Int64()

				planetInfos.Player.IsBandit = strings.Contains(player.Rank.RankClass, "bandit")
				planetInfos.Player.IsStarlord = strings.Contains(player.Rank.RankClass, "starlord")

				// Alliance
				if player.AllianceID.Int64() > 0 {
					planetInfos.Alliance = new(ogame.AllianceInfos)
					planetInfos.Alliance.ID = player.AllianceID.Int64()
					planetInfos.Alliance.Name = player.AllianceName
					planetInfos.Alliance.Tag = player.AllianceTag
					planetInfos.Alliance.Rank = player.PositionAlliance.Int64()
					if player.IsAllianceMember {
						planetInfos.Alliance.Member = 1
					}
				}

			} else if planet.PlanetType.Int64() == ogame.MoonType.Int64() {
				// Moon
				planetInfos.Moon = new(ogame.MoonInfos)

				planetInfos.Moon.ID = planet.PlanetID.Int64()
				planetInfos.Moon.Name = planet.PlanetName
				planetInfos.Moon.Diameter = int64(planet.Size)

				planetInfos.Moon.Activity = planet.Activity.Minutes

			} else if planet.PlanetType.Int64() == ogame.DebrisType.Int64() {
				// Debris Field
				planetInfos.Debris.Metal = planet.Resources.Metal
				planetInfos.Debris.Crystal = planet.Resources.Crystal
				planetInfos.Debris.Deuterium = planet.Resources.Deuterium
				planetInfos.Debris.RecyclersNeeded = planet.RequiredShips.Int64()
			}
		}
		res.SetPlanet(i, planetInfos)
	}

	return res, nil
}
