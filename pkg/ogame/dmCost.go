package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
	"strconv"
)

// DMCost ...
type DMCost struct {
	Cost                int64
	CanBuy              bool  // Either or not we have enough DM
	Complete            bool  // false means we will halve the time, true will complete
	OGameID             ID    // What we are going to build
	Nbr                 int64 // Either the amount of ships/defences or the building/research level
	BuyAndActivateToken string
	Token               string
}

// String ...
func (d DMCost) String() string {
	return "\n" +
		"               Cost: " + utils.FI64(d.Cost) + "\n" +
		"             CanBuy: " + strconv.FormatBool(d.CanBuy) + "\n" +
		"           Complete: " + strconv.FormatBool(d.Complete) + "\n" +
		"            OGameID: " + utils.FI64(d.OGameID) + "\n" +
		"                Nbr: " + utils.FI64(d.Nbr) + "\n" +
		"BuyAndActivateToken: " + d.BuyAndActivateToken + "\n" +
		"              Token: " + d.Token
}

// DMCosts ...
type DMCosts struct {
	Buildings DMCost
	Research  DMCost
	Shipyard  DMCost
}

// String ...
func (d DMCosts) String() string {
	return "\n" +
		"Buildings:" + d.Buildings.String() + "\n" +
		"Research:" + d.Research.String() + "\n" +
		"Shipyard:" + d.Shipyard.String()
}
