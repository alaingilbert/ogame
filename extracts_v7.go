package ogame

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func extractFacilitiesFromDocV7(doc *goquery.Document) (Facilities, error) {
	getVal := func(name string) int {
		val, _ := strconv.Atoi(doc.Find("span."+name+" span").First().AttrOr("data-value", "0"))
		return val
	}
	res := Facilities{}
	res.RoboticsFactory = getVal("roboticsFactory")
	res.Shipyard = getVal("shipyard")
	res.ResearchLab = getVal("researchLaboratory")
	res.AllianceDepot = getVal("allianceDepot")
	res.MissileSilo = getVal("missileSilo")
	res.NaniteFactory = getVal("naniteFactory")
	res.Terraformer = getVal("terraformer")
	res.SpaceDock = getVal("repairDock")
	res.LunarBase = getVal("lunarBase")         // TODO: ensure name is correct
	res.SensorPhalanx = getVal("sensorPhalanx") // TODO: ensure name is correct
	res.JumpGate = getVal("jumpGate")           // TODO: ensure name is correct
	return res, nil
}
