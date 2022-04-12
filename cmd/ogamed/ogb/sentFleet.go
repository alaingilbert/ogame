package ogb

import "github.com/faunX/ogame"

type SentFleet struct {
	OriginID          ogame.CelestialID
	OriginCoords      ogame.Coordinate
	DestinationCoords ogame.Coordinate
	Mission           ogame.MissionID
	Speed             ogame.Speed
	HoldingTime       int64
	Ships             ogame.ShipsInfos
	Resources         ogame.Resources
	Token             string
}
