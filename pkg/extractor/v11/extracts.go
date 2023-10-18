package v11

import (
	"errors"
	"regexp"
	"bytes"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

func extractResourceSettingsFromPage(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID := v6.ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ResourceSettings{}, "", ogame.ErrInvalidPlanetID
	}
	vals := make([]int64, 0)
	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		_, selectedExists := s.Attr("selected")
		if selectedExists {
			a, _ := s.Attr("value")
			val := utils.DoParseI64(a)
			vals = append(vals, val)
		}
	})
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

	getToken := func(pageHTML []byte) (string, error) {
		m := regexp.MustCompile(`var token = "([^"]+)"`).FindSubmatch(pageHTML)
		if len(m) != 2 {
			return "", errors.New("unable to find token")
		}
		return string(m[1]), nil
	}
	token, _ := getToken(pageHTML)

	return res, token, nil
}