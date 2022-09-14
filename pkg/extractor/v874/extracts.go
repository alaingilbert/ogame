package v874

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
)

func extractBuffActivationFromDoc(doc *goquery.Document) (token string, items []ogame.Item, err error) {
	scriptTxt := doc.Find("script").Text()
	r := regexp.MustCompile(`token = "([^"]+)"`)
	m := r.FindStringSubmatch(scriptTxt)
	if len(m) != 2 {
		err = errors.New("failed to find activate token")
		return
	}
	token = m[1]
	r = regexp.MustCompile(`items_inventory = ({[^\n]+});\n`)
	m = r.FindStringSubmatch(scriptTxt)
	if len(m) != 2 {
		err = errors.New("failed to find items inventory")
		return
	}
	var inventoryMap map[string]ogame.Item
	if err = json.Unmarshal([]byte(m[1]), &inventoryMap); err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range inventoryMap {
		items = append(items, item)
	}
	return
}

func extractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error) {
	s := doc.Find("div.js_import_price")
	if s.Size() == 0 {
		err = errors.New("failed to extract offer of the day price")
		return
	}
	price = utils.ParseInt(s.Text())
	script := doc.Find("script").Text()
	m := regexp.MustCompile(`var token\s?=\s?"([^"]*)";`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day import token")
		return
	}
	importToken = string(m[1])
	m = regexp.MustCompile(`var planetResources\s?=\s?({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw planet resources")
		return
	}
	if err = json.Unmarshal(m[1], &planetResources); err != nil {
		return
	}
	m = regexp.MustCompile(`var multiplier\s?=\s?({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw multiplier")
		return
	}
	if err = json.Unmarshal(m[1], &multiplier); err != nil {
		return
	}
	return
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
		endAtApprox := doc.Find("p.auction_info b").Text()
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
