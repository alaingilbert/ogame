package v11_9_0

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"regexp"
)

func extractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	res := make([]ogame.Quantifiable, 0)
	active := doc.Find("table.construction")
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []ogame.Quantifiable{}, nil
	}
	idInt := utils.DoParseI64(m[1])
	activeID := ogame.ID(idInt)
	activeNbr := utils.DoParseI64(active.Find("div.shipSumCount").Text())
	res = append(res, ogame.Quantifiable{ID: activeID, Nbr: activeNbr})
	doc.Find("table.queue td").Each(func(i int, s *goquery.Selection) {
		link := s.Find("img")
		alt := link.AttrOr("alt", "")
		var itemID ogame.ID
		if id := ogame.DefenceName2ID(alt); id.IsValid() {
			itemID = id
		} else if id := ogame.ShipName2ID(alt); id.IsValid() {
			itemID = id
		}
		if itemID.IsValid() {
			itemNbr := utils.ParseInt(s.Text())
			res = append(res, ogame.Quantifiable{ID: ogame.ID(itemID), Nbr: itemNbr})
		}
	})
	return res, nil
}
