package ogame

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
		"RoboticsFactory: " + FI64(f.RoboticsFactory) + "\n" +
		"       Shipyard: " + FI64(f.Shipyard) + "\n" +
		"   Research Lab: " + FI64(f.ResearchLab) + "\n" +
		" Alliance Depot: " + FI64(f.AllianceDepot) + "\n" +
		"   Missile Silo: " + FI64(f.MissileSilo) + "\n" +
		" Nanite Factory: " + FI64(f.NaniteFactory) + "\n" +
		"    Terraformer: " + FI64(f.Terraformer) + "\n" +
		"     Space Dock: " + FI64(f.SpaceDock) + "\n" +
		"     Lunar Base: " + FI64(f.LunarBase) + "\n" +
		" Sensor Phalanx: " + FI64(f.SensorPhalanx) + "\n" +
		"      Jump Gate: " + FI64(f.JumpGate)
}
