package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
	"strconv"
)

// Auction ...
type Auction struct {
	Ref                 string
	HasFinished         bool
	Endtime             int64
	NumBids             int64
	CurrentBid          int64
	AlreadyBid          int64
	MinimumBid          int64
	DeficitBid          int64
	HighestBidder       string
	HighestBidderUserID int64
	CurrentItem         string
	CurrentItemLong     string
	Inventory           int64
	Token               string
	ResourceMultiplier  struct {
		Metal     float64
		Crystal   float64
		Deuterium float64
		Honor     int64
	}
	Resources map[string]any
}

// String ...
func (a Auction) String() string {
	return "" +
		"           Ref: " + a.Ref + "\n" +
		"  Has finished: " + strconv.FormatBool(a.HasFinished) + "\n" +
		"      End time: " + utils.FI64(a.Endtime) + "\n" +
		"      Num bids: " + utils.FI64(a.NumBids) + "\n" +
		"   Minimum bid: " + utils.FI64(a.MinimumBid) + "\n" +
		"Highest bidder: " + a.HighestBidder + " (" + utils.FI64(a.HighestBidderUserID) + ")" + "\n" +
		"  Current item: " + a.CurrentItem + " (" + a.CurrentItemLong + ")" + "\n" +
		"     Inventory: " + utils.FI64(a.Inventory) + "\n" +
		""
}
