package v11_13_0

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"regexp"
)

func ExtractConstructions(pageHTML []byte, clock clockwork.Clock) (buildingID ogame.ID, buildingCountdown int64,
	researchID ogame.ID, researchCountdown int64,
	lfBuildingID ogame.ID, lfBuildingCountdown int64,
	lfResearchID ogame.ID, lfResearchCountdown int64) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	buildingDataEnd := utils.DoParseI64(doc.Find("time.buildingCountdown").AttrOr("data-end", "0"))
	if buildingDataEnd > 0 {
		buildingCountdown = buildingDataEnd - clock.Now().Unix()
		buildingIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelbuilding\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ogame.ID(buildingIDInt)
	}
	researchDataEnd := utils.DoParseI64(doc.Find("time.researchCountdown").AttrOr("data-end", "0"))
	if researchDataEnd > 0 {
		researchCountdown = researchDataEnd - clock.Now().Unix()
		researchIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelresearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ogame.ID(researchIDInt)
	}
	lfBuildingDataEnd := utils.DoParseI64(doc.Find("time.lfbuildingCountdown").AttrOr("data-end", "0"))
	if lfBuildingDataEnd > 0 {
		lfBuildingCountdown = lfBuildingDataEnd - clock.Now().Unix()
		lfBuildingIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancellfbuilding\((\d+),`).FindSubmatch(pageHTML)[1])
		lfBuildingID = ogame.ID(lfBuildingIDInt)
	}
	lfResearchDataEnd := utils.DoParseI64(doc.Find("time.lfResearchCountdown").AttrOr("data-end", "0"))
	if lfResearchDataEnd > 0 {
		lfResearchCountdown = lfResearchDataEnd - clock.Now().Unix()
		lfResearchIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancellfresearch\((\d+),`).FindSubmatch(pageHTML)[1])
		lfResearchID = ogame.ID(lfResearchIDInt)
	}
	return
}
