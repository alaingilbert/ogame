package ogame

import "strconv"

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

// Lazy returns a function that return self
func (f Facilities) Lazy() LazyFacilities {
	return func() Facilities { return f }
}

// ByID gets the facility level by facility id
func (f Facilities) ByID(id ID) int64 {
	if id == RoboticsFactory.ID {
		return f.RoboticsFactory
	} else if id == Shipyard.ID {
		return f.Shipyard
	} else if id == ResearchLab.ID {
		return f.ResearchLab
	} else if id == AllianceDepot.ID {
		return f.AllianceDepot
	} else if id == MissileSilo.ID {
		return f.MissileSilo
	} else if id == NaniteFactory.ID {
		return f.NaniteFactory
	} else if id == Terraformer.ID {
		return f.Terraformer
	} else if id == SpaceDock.ID {
		return f.SpaceDock
	} else if id == LunarBase.ID {
		return f.LunarBase
	} else if id == SensorPhalanx.ID {
		return f.SensorPhalanx
	} else if id == JumpGate.ID {
		return f.JumpGate
	}
	return 0
}

func (f Facilities) String() string {
	return "\n" +
		"RoboticsFactory: " + strconv.FormatInt(f.RoboticsFactory, 10) + "\n" +
		"       Shipyard: " + strconv.FormatInt(f.Shipyard, 10) + "\n" +
		"   Research Lab: " + strconv.FormatInt(f.ResearchLab, 10) + "\n" +
		" Alliance Depot: " + strconv.FormatInt(f.AllianceDepot, 10) + "\n" +
		"   Missile Silo: " + strconv.FormatInt(f.MissileSilo, 10) + "\n" +
		" Nanite Factory: " + strconv.FormatInt(f.NaniteFactory, 10) + "\n" +
		"    Terraformer: " + strconv.FormatInt(f.Terraformer, 10) + "\n" +
		"     Space Dock: " + strconv.FormatInt(f.SpaceDock, 10) + "\n" +
		"     Lunar Base: " + strconv.FormatInt(f.LunarBase, 10) + "\n" +
		" Sensor Phalanx: " + strconv.FormatInt(f.SensorPhalanx, 10) + "\n" +
		"      Jump Gate: " + strconv.FormatInt(f.JumpGate, 10)
}
