package ogame

import "strconv"

// Facilities represent a planet facilities information
type Facilities struct {
	RoboticsFactory int
	Shipyard        int
	ResearchLab     int
	AllianceDepot   int
	MissileSilo     int
	NaniteFactory   int
	Terraformer     int
	SpaceDock       int
}

// ByID gets the facility level by facility id
func (f Facilities) ByID(id ID) int {
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
	}
	return 0
}

func (f Facilities) String() string {
	return "\n" +
		"RoboticsFactory: " + strconv.Itoa(f.RoboticsFactory) + "\n" +
		"       Shipyard: " + strconv.Itoa(f.Shipyard) + "\n" +
		"   Research Lab: " + strconv.Itoa(f.ResearchLab) + "\n" +
		" Alliance Depot: " + strconv.Itoa(f.AllianceDepot) + "\n" +
		"   Missile Silo: " + strconv.Itoa(f.MissileSilo) + "\n" +
		" Nanite Factory: " + strconv.Itoa(f.NaniteFactory) + "\n" +
		"    Terraformer: " + strconv.Itoa(f.Terraformer) + "\n" +
		"     Space Dock: " + strconv.Itoa(f.SpaceDock)
}
