package v11_13_0

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ExtractConstructions(pageHTML []byte, clock clockwork.Clock) (out ogame.Constructions, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return
	}
	data := [][]string{
		{"buildingCountdown", "cancelbuilding"},
		{"researchCountdown", "cancelresearch"},
		{"lfbuildingCountdown", "cancellfbuilding"},
		{"lfResearchCountdown", "cancellfresearch"},
	}
	constructionRows := make([]ogame.Construction, 4)
	for i, d := range data {
		s := doc.Find("time." + d[0])
		parent := s.Parent().Parent().Parent()
		buildingDataEnd := utils.DoParseI64(s.AttrOr("data-end", "0"))
		if buildingDataEnd > 0 {
			countdown := time.Duration(buildingDataEnd-clock.Now().Unix()) * time.Second
			id := ogame.ID(utils.ToInt(regexp.MustCompile(`onclick="` + d[1] + `\((\d+),`).FindSubmatch(pageHTML)[1]))
			level := utils.DoParseI64(regexp.MustCompile(`(\d+)`).FindStringSubmatch(parent.Find("span.level").Text())[1])
			constructionRows[i] = ogame.Construction{ID: id, Countdown: countdown, Level: level}
		}
	}
	out.Building = constructionRows[0]
	out.Research = constructionRows[1]
	out.LfBuilding = constructionRows[2]
	out.LfResearch = constructionRows[3]
	return
}

func extractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	var shipSumCountdown int64
	shipSumCountdownMatch := regexp.MustCompile(`CountdownTimer\('shipyardCountdown', (\d+),`).FindSubmatch(pageHTML)
	if len(shipSumCountdownMatch) > 0 {
		shipSumCountdown = int64(utils.ToInt(shipSumCountdownMatch[1]))
	}
	return shipSumCountdown
}

func extractFleetsFromDoc(doc *goquery.Document, location *time.Location, lifeformEnabled bool) (res []ogame.Fleet, err error) {
	res = make([]ogame.Fleet, 0)
	script := doc.Find("body script").Text()
	for _, s := range doc.Find("div.fleetDetails").EachIter() {
		originText := s.Find("span.originCoords a").Text()
		origin := v6.ExtractCoord(originText)
		origin.Type = ogame.PlanetType
		if s.Find("span.originPlanet figure").HasClass("moon") {
			origin.Type = ogame.MoonType
		}

		destText := s.Find("span.destinationCoords a").Text()
		dest := v6.ExtractCoord(destText)
		dest.Type = ogame.PlanetType
		if s.Find("span.destinationPlanet figure").HasClass("moon") {
			dest.Type = ogame.MoonType
		} else if s.Find("span.destinationPlanet figure").HasClass("tf") {
			dest.Type = ogame.DebrisType
		}

		id := utils.DoParseI64(s.Find("a.openCloseDetails").AttrOr("data-mission-id", "0"))

		timerID := strings.TrimPrefix(s.Find("span.timer").AttrOr("id", ""), "timer_")
		m := regexp.MustCompile(`SimpleCountdownTimer\(\s*"#timer_` + timerID + `",\s*(\d+),`).FindStringSubmatch(script)
		var arriveIn int64
		if len(m) == 2 {
			arriveIn = utils.DoParseI64(m[1])
		}

		timerNextID := s.Find("span.nextTimer").AttrOr("id", "")
		m = regexp.MustCompile(`getElementByIdWithCache\("` + timerNextID + `"\),\s*(\d+)\s*\);`).FindStringSubmatch(script)
		var backIn int64
		if len(m) == 2 {
			backIn = utils.DoParseI64(m[1])
		}

		missionType := utils.DoParseI64(s.AttrOr("data-mission-type", ""))
		returnFlight, _ := strconv.ParseBool(s.AttrOr("data-return-flight", ""))
		inDeepSpace := s.Find("span.fleetDetailButton a").HasClass("fleet_icon_forward_end")
		arrivalTime := utils.DoParseI64(s.AttrOr("data-arrival-time", ""))
		endTime := utils.DoParseI64(s.Find("a.openCloseDetails").AttrOr("data-end-time", ""))

		trs := s.Find("table.fleetinfo tr")
		shipment := ogame.Resources{}
		metalTrOffset := 3
		crystalTrOffset := 2
		DeuteriumTrOffset := 1
		if lifeformEnabled {
			metalTrOffset = 4
			crystalTrOffset = 3
			DeuteriumTrOffset = 2
		}
		shipment.Metal = utils.ParseInt(trs.Eq(trs.Size() - metalTrOffset).Find("td").Eq(1).Text())
		shipment.Crystal = utils.ParseInt(trs.Eq(trs.Size() - crystalTrOffset).Find("td").Eq(1).Text())
		shipment.Deuterium = utils.ParseInt(trs.Eq(trs.Size() - DeuteriumTrOffset).Find("td").Eq(1).Text())

		fedAttackHref := s.Find("span.fedAttack a").AttrOr("href", "")
		fedAttackURL, err := url.Parse(fedAttackHref)
		if err != nil {
			return nil, err
		}
		fedAttackQuery := fedAttackURL.Query()
		targetPlanetID := utils.DoParseI64(fedAttackQuery.Get("target"))
		unionID := utils.DoParseI64(fedAttackQuery.Get("union"))

		fleet := ogame.MakeFleet()
		fleet.ID = ogame.FleetID(id)
		fleet.Origin = origin
		fleet.Destination = dest
		fleet.Mission = ogame.MissionID(missionType)
		fleet.ReturnFlight = returnFlight
		fleet.InDeepSpace = inDeepSpace
		fleet.Resources = shipment
		fleet.TargetPlanetID = targetPlanetID
		fleet.UnionID = unionID
		fleet.ArrivalTime = time.Unix(endTime, 0)
		fleet.BackTime = time.Unix(arrivalTime, 0)

		var startTimeString string
		var startTimeStringExists bool
		if !returnFlight {
			fleet.ArriveIn = arriveIn
			fleet.BackIn = backIn
			startTimeString, startTimeStringExists = s.Find("div.origin img").Attr("title")
		} else {
			fleet.ArriveIn = -1
			fleet.BackIn = arriveIn
			startTimeString, startTimeStringExists = s.Find("div.destination img").Attr("title")
		}

		var startTime time.Time
		if startTimeStringExists {
			startTimeArray := strings.Split(startTimeString, ":| ")
			if len(startTimeArray) == 2 {
				startTime, _ = time.ParseInLocation("02.01.2006<br>15:04:05", startTimeArray[1], location)
			}
		}
		fleet.StartTime = startTime.Local()

		for i := 1; i < trs.Size()-5; i++ {
			tds := trs.Eq(i).Find("td")
			name := strings.ToLower(strings.Trim(strings.TrimSpace(tds.Eq(0).Text()), ":"))
			qty := utils.ParseInt(tds.Eq(1).Text())
			shipID := ogame.ShipName2ID(name)
			fleet.Ships.Set(shipID, qty)
		}

		res = append(res, fleet)
	}
	return
}
