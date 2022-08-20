package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
)

// LazyFacilities ...
type LazyFacilities func() Facilities

// Facilities represent a planet facilities information
type Facilities struct {
	RoboticsFactory int64 // 14
	Shipyard        int64 // 21
	ResearchLab     int64 // 31
	AllianceDepot   int64 // 34
	MissileSilo     int64 // 44
	NaniteFactory   int64 // 15
	Terraformer     int64 // 33
	SpaceDock       int64 // 36
	LunarBase       int64 // 41
	SensorPhalanx   int64 // 42
	JumpGate        int64 // 43
}

func (f Facilities) GetNaniteFactory() int64   { return f.NaniteFactory }
func (f Facilities) GetRoboticsFactory() int64 { return f.RoboticsFactory }
func (f Facilities) GetResearchLab() int64     { return f.ResearchLab }
func (f Facilities) GetShipyard() int64        { return f.Shipyard }

// Lazy returns a function that return self
func (f Facilities) Lazy() LazyFacilities {
	return func() Facilities { return f }
}

// ByID gets the facility level by facility id
func (f Facilities) ByID(id ID) int64 {
	if id == RoboticsFactoryID {
		return f.RoboticsFactory
	} else if id == ShipyardID {
		return f.Shipyard
	} else if id == ResearchLabID {
		return f.ResearchLab
	} else if id == AllianceDepotID {
		return f.AllianceDepot
	} else if id == MissileSiloID {
		return f.MissileSilo
	} else if id == NaniteFactoryID {
		return f.NaniteFactory
	} else if id == TerraformerID {
		return f.Terraformer
	} else if id == SpaceDockID {
		return f.SpaceDock
	} else if id == LunarBaseID {
		return f.LunarBase
	} else if id == SensorPhalanxID {
		return f.SensorPhalanx
	} else if id == JumpGateID {
		return f.JumpGate
	}
	return 0
}

func (f Facilities) String() string {
	return "\n" +
		"RoboticsFactory: " + utils.FI64(f.RoboticsFactory) + "\n" +
		"       Shipyard: " + utils.FI64(f.Shipyard) + "\n" +
		"   Research Lab: " + utils.FI64(f.ResearchLab) + "\n" +
		" Alliance Depot: " + utils.FI64(f.AllianceDepot) + "\n" +
		"   Missile Silo: " + utils.FI64(f.MissileSilo) + "\n" +
		" Nanite Factory: " + utils.FI64(f.NaniteFactory) + "\n" +
		"    Terraformer: " + utils.FI64(f.Terraformer) + "\n" +
		"     Space Dock: " + utils.FI64(f.SpaceDock) + "\n" +
		"     Lunar Base: " + utils.FI64(f.LunarBase) + "\n" +
		" Sensor Phalanx: " + utils.FI64(f.SensorPhalanx) + "\n" +
		"      Jump Gate: " + utils.FI64(f.JumpGate)
}
