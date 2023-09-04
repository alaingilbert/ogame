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
		Galaxy            int64         `json:"galaxy"`
		System            int64         `json:"system"`
		DeuteriumInDebris bool          `json:"deuteriumInDebris"`
		Content           []ContentJson `json:"galaxyContent"`
	}

	ContentJson struct {
		Galaxy   int64       `json:"galaxy"`
		System   int64       `json:"system"`
		Position int64       `json:"position"`
		Planets  PlanetsList `json:"planets"`
		Player   PlayerJson  `json:"player"`
		Filters  string      `json:"positionFilters"`
	}

	PlanetsList []PlanetJson

	PlanetJson struct {
		PlayerID      int64         `json:"playerId"`
		PlanetID      int64         `json:"planetId"`
		PlanetName    string        `json:"planetName"`
		Image         string        `json:"imageInformation"`
		PlanetType    int64         `json:"planetType"`
		IsDestroyed   bool          `json:"isDestroyed"`
		Size          int64         `json:"size,string"`
		RequiredShips int64         `json:"requiredShips"`
		Resources     ResourcesJson `json:"resources"`
		Activity      ActivityJson  `json:"activity"`
	}

	ResourcesJson struct {
		Metal struct {
			Amount int64 `json:"amount,string"`
		} `json:"metal"`
		Crystal struct {
			Amount int64 `json:"amount,string"`
		} `json:"crystal"`
		Deuterium struct {
			Amount int64 `json:"amount,string"`
		} `json:"deuterium"`
	}

	ActivityJson struct {
		Minutes int64
	}

	PlayerJson struct {
		PlayerID          int64    `json:"playerId"`
		PlayerName        string   `json:"playerName"`
		AllianceID        int64    `json:"allianceId"`
		AllianceName      string   `json:"allianceName"`
		AllianceTag       string   `json:"allianceTag"`
		IsAllianceMember  bool     `json:"isAllianceMember"`
		PositionPlayer    int64    `json:"highscorePositionPlayer"`
		PositionAlliance  int64    `json:"highscorePositionAlliance,string"`
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
)

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

func (a *ActivityJson) UnmarshalJSON(d []byte) error {

	var min int64
	var tmp map[string]interface{}
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
	res.SetGalaxy(tmp.System.Galaxy)
	res.SetSystem(tmp.System.System)

	for i, pos := range tmp.System.Content {

		if pos.Position == 16 {
			res.ExpeditionDebris.Metal = pos.Planets[0].Resources.Metal.Amount
			res.ExpeditionDebris.Crystal = pos.Planets[0].Resources.Crystal.Amount
			res.ExpeditionDebris.Deuterium = pos.Planets[0].Resources.Deuterium.Amount
			res.ExpeditionDebris.PathfindersNeeded = pos.Planets[0].RequiredShips
			continue
		}
		// TODO: manage asteroids and events in P17
		if pos.Position == 17 {
			continue
		}
		if len(pos.Planets) == 0 {
			continue
		}
		player := pos.Player
		planetInfos := new(ogame.PlanetInfos)

		for _, planet := range pos.Planets {

			// Planet
			if planet.PlanetType == ogame.PlanetType.Int64() {
				// Generic infos
				planetInfos.ID = planet.PlanetID
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

				planetInfos.Coordinate.Galaxy = pos.Galaxy
				planetInfos.Coordinate.System = pos.System
				planetInfos.Coordinate.Position = pos.Position
				planetInfos.Coordinate.Type = ogame.PlanetType

				planetInfos.Date = time.Now()

				// Player
				planetInfos.Player.ID = player.PlayerID
				planetInfos.Player.Name = player.PlayerName
				planetInfos.Player.Rank = player.PositionPlayer

				planetInfos.Player.IsBandit = strings.Contains(player.Rank.RankClass, "bandit")
				planetInfos.Player.IsStarlord = strings.Contains(player.Rank.RankClass, "starlord")

				// Alliance
				if player.AllianceID > 0 {
					planetInfos.Alliance = new(ogame.AllianceInfos)
					planetInfos.Alliance.ID = player.AllianceID
					planetInfos.Alliance.Name = player.AllianceName
					planetInfos.Alliance.Tag = player.AllianceTag
					planetInfos.Alliance.Rank = player.PositionAlliance
					if player.IsAllianceMember {
						planetInfos.Alliance.Member = 1
					}
				}

			} else if planet.PlanetType == ogame.MoonType.Int64() {
				// Moon
				planetInfos.Moon = new(ogame.MoonInfos)

				planetInfos.Moon.ID = planet.PlanetID
				planetInfos.Moon.Name = planet.PlanetName
				planetInfos.Moon.Diameter = planet.Size

				planetInfos.Moon.Activity = planet.Activity.Minutes

			} else if planet.PlanetType == ogame.DebrisType.Int64() {
				// Debris Field
				planetInfos.Debris.Metal = planet.Resources.Metal.Amount
				planetInfos.Debris.Crystal = planet.Resources.Crystal.Amount
				planetInfos.Debris.Deuterium = planet.Resources.Deuterium.Amount
				planetInfos.Debris.RecyclersNeeded = planet.RequiredShips
			}
		}
		res.SetPlanet(i, planetInfos)
	}

	return res, nil
}
