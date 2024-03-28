package v11_9_0

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
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

func extractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error) {
	msgs := make([]ogame.CombatReportSummary, 0)
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				report := ogame.CombatReportSummary{ID: id}
				report.Destination = v6.ExtractCoord(s.Find("div.msg_head a").Text())
				if s.Find("div.msg_head figure").HasClass("planet") {
					report.Destination.Type = ogame.PlanetType
				} else if s.Find("div.msg_head figure").HasClass("moon") {
					report.Destination.Type = ogame.MoonType
				} else {
					report.Destination.Type = ogame.PlanetType
				}
				apiKeyTitle := s.Find("span.icon_apikey").AttrOr("title", "")
				m := regexp.MustCompile(`'(cr-[^']+)'`).FindStringSubmatch(apiKeyTitle)
				if len(m) == 2 {
					report.APIKey = m[1]
				}
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m = regexp.MustCompile(`([\d.,]+)<br/>[^\d]*([\d.,]+)<br/>[^\d]*([\d.,]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = utils.ParseInt(m[1])
					report.Crystal = utils.ParseInt(m[2])
					report.Deuterium = utils.ParseInt(m[3])
				}
				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = utils.ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.,]+[^\d]*([\d.,]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = utils.ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("span.msg_date").Text())
				report.CreatedAt = msgDate

				link := s.Find("message-footer.msg_actions button.msgAttackBtn").AttrOr("onclick", "")
				m = regexp.MustCompile(`page=ingame&component=fleetdispatch&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
				if len(m) != 5 {
					return
				}
				galaxy := utils.DoParseI64(m[1])
				system := utils.DoParseI64(m[2])
				position := utils.DoParseI64(m[3])
				planetType := utils.DoParseI64(m[4])
				report.Origin = &ogame.Coordinate{Galaxy: galaxy, System: system, Position: position, Type: ogame.CelestialType(planetType)}
				if report.Origin.Equal(report.Destination) {
					report.Origin = nil
				}

				msgs = append(msgs, report)
			}
		}
	})
	return msgs, nbPage, nil
}

// extractAuctionFromDoc extract auction information from page "traderAuctioneer"
func extractAuctionFromDoc(doc *goquery.Document) (ogame.Auction, error) {
	auction := ogame.Auction{}
	auction.HasFinished = false

	// Detect if Auction has already finished
	nextAuction := doc.Find("#nextAuction")
	if nextAuction.Size() > 0 {
		// Find time until next auction starts
		auction.Endtime = utils.DoParseI64(nextAuction.Text())
		auction.HasFinished = true
	} else {
		endAtApprox := doc.Find("p.auction_info span").Text()
		m := regexp.MustCompile(`[^\d]+(\d+).*`).FindStringSubmatch(endAtApprox)
		if len(m) != 2 {
			return ogame.Auction{}, errors.New("failed to find end time approx")
		}
		endTimeMinutes, err := utils.ParseI64(m[1])
		if err != nil {
			return ogame.Auction{}, errors.New("invalid end time approx: " + err.Error())
		}
		auction.Endtime = endTimeMinutes * 60
	}

	auction.HighestBidder = strings.TrimSpace(doc.Find("a.currentPlayer").Text())
	auction.HighestBidderUserID = utils.DoParseI64(doc.Find("a.currentPlayer").AttrOr("data-player-id", ""))
	auction.NumBids = utils.DoParseI64(doc.Find("div.numberOfBids").Text())
	auction.CurrentBid = utils.ParseInt(doc.Find("div.currentSum").Text())
	auction.Inventory = utils.DoParseI64(doc.Find("span.level.amount").Text())
	auction.CurrentItem = strings.ToLower(doc.Find("img").First().AttrOr("alt", ""))
	auction.CurrentItemLong = strings.ToLower(doc.Find("div.image_140px").First().Find("a").First().AttrOr("title", ""))
	multiplierRegex := regexp.MustCompile(`multiplier\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(multiplierRegex) != 2 {
		return ogame.Auction{}, errors.New("failed to find auction multiplier")
	}
	if err := json.Unmarshal([]byte(multiplierRegex[1]), &auction.ResourceMultiplier); err != nil {
		return ogame.Auction{}, errors.New("failed to json parse auction multiplier: " + err.Error())
	}

	// Find auctioneer token
	tokenRegex := regexp.MustCompile(`token\s?=\s?"([^"]+)";`).FindStringSubmatch(doc.Text())
	if len(tokenRegex) != 2 {
		return ogame.Auction{}, errors.New("failed to find auctioneer token")
	}
	auction.Token = tokenRegex[1]

	// Find Planet / Moon resources JSON
	planetMoonResources := regexp.MustCompile(`planetResources\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(planetMoonResources) != 2 {
		return ogame.Auction{}, errors.New("failed to find planetResources")
	}
	if err := json.Unmarshal([]byte(planetMoonResources[1]), &auction.Resources); err != nil {
		return ogame.Auction{}, errors.New("failed to json unmarshal planetResources: " + err.Error())
	}

	// Find already-bid
	m := regexp.MustCompile(`var playerBid\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(m) != 2 {
		return ogame.Auction{}, errors.New("failed to get playerBid")
	}
	var alreadyBid int64
	if m[1] != "false" {
		alreadyBid = utils.DoParseI64(m[1])
	}
	auction.AlreadyBid = alreadyBid

	// Find min-bid
	auction.MinimumBid = utils.ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_price").Text())

	// Find deficit-bid
	auction.DeficitBid = utils.ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_deficit").Text())

	// Note: Don't just bid the min-bid amount. It will keep doubling the total bid and grow exponentially...
	// DeficitBid is 1000 when another player has outbid you or if nobody has bid yet.
	// DeficitBid seems to be filled by Javascript in the browser. We're parsing it anyway. Correct Bid calculation would be:
	// bid = max(auction.DeficitBid, auction.MinimumBid - auction.AlreadyBid)

	return auction, nil
}
