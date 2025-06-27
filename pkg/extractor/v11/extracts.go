package v11

import (
	"bytes"
	"errors"
	v104 "github.com/alaingilbert/ogame/pkg/extractor/v104"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func extractResourceSettingsFromPage(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ResourceSettings{}, "", err
	}
	bodyID := v6.ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ResourceSettings{}, "", ogame.ErrInvalidPlanetID
	}
	vals := make([]int64, 0)
	for _, s := range doc.Find("option").EachIter() {
		if _, selectedExists := s.Attr("selected"); selectedExists {
			val := utils.DoParseI64(s.AttrOr("value", ""))
			vals = append(vals, val)
		}
	}
	if len(vals) != 7 {
		return ogame.ResourceSettings{}, "", errors.New("failed to find all resource settings")
	}

	res := ogame.ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]
	res.Crawler = vals[6]

	token, _ := v104.ExtractToken(pageHTML)

	return res, token, nil
}

func extractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return ExtractCancelInfos(pageHTML, "cancelbuilding", 0)
}

func extractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return ExtractCancelInfos(pageHTML, "cancelresearch", 2)
}

func extractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error) {
	return ExtractCancelInfos(pageHTML, "cancellfbuilding", 1)
}

func ExtractCancelInfos(pageHTML []byte, fnName string, tableIdx int) (token string, id, listID int64, err error) {
	r1 := regexp.MustCompile(`window\.token = '([^']+)'`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return "", 0, 0, err
	}
	t := doc.Find("table.construction").Eq(tableIdx)
	a := t.Find("a").First().AttrOr("onclick", "")
	r := regexp.MustCompile(fnName + `\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find id/listid")
	}
	id = utils.DoParseI64(m[1])
	listID = utils.DoParseI64(m[2])
	return
}

func extractCharacterClassFromDoc(doc *goquery.Document) ogame.CharacterClass {
	characterClassDiv := doc.Find("div#characterclass a div")
	characterClass := ogame.NoClass
	if characterClassDiv.HasClass("miner") {
		characterClass = ogame.Collector
	} else if characterClassDiv.HasClass("warrior") {
		characterClass = ogame.General
	} else if characterClassDiv.HasClass("explorer") {
		characterClass = ogame.Discoverer
	}
	return characterClass
}

func extractLifeformTypeFromDoc(doc *goquery.Document) ogame.LifeformType {
	lfDiv := doc.Find("div#lifeform div.lifeform-item-icon")
	if lfDiv.HasClass("lifeform1") {
		return ogame.Humans
	} else if lfDiv.HasClass("lifeform2") {
		return ogame.Rocktal
	} else if lfDiv.HasClass("lifeform3") {
		return ogame.Mechas
	} else if lfDiv.HasClass("lifeform4") {
		return ogame.Kaelesh
	}
	return ogame.NoneLfType
}

func extractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64, error) {
	m := regexp.MustCompile(`\$\("#cooldown"\), (\d+),`).FindSubmatch(pageHTML)
	ships := ogame.ShipsInfos{}
	var destinations []ogame.MoonID
	if len(m) > 0 {
		waitTime := int64(utils.ToInt(m[1]))
		return ships, "", destinations, waitTime, nil
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ships, "", destinations, 0, err
	}
	for _, s := range ogame.Ships {
		ships.Set(s.GetID(), utils.ParseInt(doc.Find("input#ship_"+utils.FI64(s.GetID())).AttrOr("rel", "0")))
	}
	token := doc.Find("input[name=token]").AttrOr("value", "")

	for _, s := range doc.Find("select[name=targetSpaceObjectId] option").EachIter() {
		moonID := utils.ParseInt(s.AttrOr("value", "0"))
		if moonID > 0 {
			destinations = append(destinations, ogame.MoonID(moonID))
		}
	}

	return ships, token, destinations, 0, nil
}

func extractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences {
	prefs := v6.ExtractPreferencesFromDoc(doc)
	prefs.Language = extractLanguageFromDoc(doc)
	return prefs
}

func extractLanguageFromDoc(doc *goquery.Document) string {
	return doc.Find("select[name=language] option[selected]").AttrOr("value", "")
}
