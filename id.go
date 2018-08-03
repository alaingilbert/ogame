package ogame

import "strconv"

// ID ...
type ID int

// String ...
func (o ID) String() string {
	res := ""
	switch o {
	case AllianceDepot.ID:
		res += "AllianceDepot"
	case RoboticsFactory.ID:
		res += "RoboticsFactory"
	case Shipyard.ID:
		res += "Shipyard"
	case ResearchLab.ID:
		res += "ResearchLab"
	case MissileSilo.ID:
		res += "MissileSilo"
	case NaniteFactory.ID:
		res += "NaniteFactory"
	case Terraformer.ID:
		res += "Terraformer"
	case SpaceDock.ID:
		res += "SpaceDock"
	case MetalMine.ID:
		res += "MetalMine"
	case CrystalMine.ID:
		res += "CrystalMine"
	case DeuteriumSynthesizer.ID:
		res += "DeuteriumSynthesizer"
	case SolarPlant.ID:
		res += "SolarPlant"
	case FusionReactor.ID:
		res += "FusionReactor"
	case MetalStorage.ID:
		res += "MetalStorage"
	case CrystalStorage.ID:
		res += "CrystalStorage"
	case DeuteriumTank.ID:
		res += "DeuteriumTank"
	case ShieldedMetalDen.ID:
		res += "ShieldedMetalDen"
	case UndergroundCrystalDen.ID:
		res += "UndergroundCrystalDen"
	case SeabedDeuteriumDen.ID:
		res += "SeabedDeuteriumDen"
	case RocketLauncher.ID:
		res += "RocketLauncher"
	case LightLaser.ID:
		res += "LightLaser"
	case HeavyLaser.ID:
		res += "HeavyLaser"
	case GaussCannon.ID:
		res += "GaussCannon"
	case IonCannon.ID:
		res += "IonCannon"
	case PlasmaTurret.ID:
		res += "PlasmaTurret"
	case SmallShieldDome.ID:
		res += "SmallShieldDome"
	case LargeShieldDome.ID:
		res += "LargeShieldDome"
	case AntiBallisticMissiles.ID:
		res += "AntiBallisticMissiles"
	case InterplanetaryMissiles.ID:
		res += "InterplanetaryMissiles"
	case SmallCargo.ID:
		res += "SmallCargo"
	case LargeCargo.ID:
		res += "LargeCargo"
	case LightFighter.ID:
		res += "LightFighter"
	case HeavyFighter.ID:
		res += "HeavyFighter"
	case Cruiser.ID:
		res += "Cruiser"
	case Battleship.ID:
		res += "Battleship"
	case ColonyShip.ID:
		res += "ColonyShip"
	case Recycler.ID:
		res += "Recycler"
	case EspionageProbe.ID:
		res += "EspionageProbe"
	case Bomber.ID:
		res += "Bomber"
	case SolarSatellite.ID:
		res += "SolarSatellite"
	case Destroyer.ID:
		res += "Destroyer"
	case Deathstar.ID:
		res += "Deathstar"
	case Battlecruiser.ID:
		res += "Battlecruiser"
	case EspionageTechnology.ID:
		res += "EspionageTechnology"
	case ComputerTechnology.ID:
		res += "ComputerTechnology"
	case WeaponsTechnology.ID:
		res += "WeaponsTechnology"
	case ShieldingTechnology.ID:
		res += "ShieldingTechnology"
	case ArmourTechnology.ID:
		res += "ArmourTechnology"
	case EnergyTechnology.ID:
		res += "EnergyTechnology"
	case HyperspaceTechnology.ID:
		res += "HyperspaceTechnology"
	case CombustionDrive.ID:
		res += "CombustionDrive"
	case ImpulseDrive.ID:
		res += "ImpulseDrive"
	case HyperspaceDrive.ID:
		res += "HyperspaceDrive"
	case LaserTechnology.ID:
		res += "LaserTechnology"
	case IonTechnology.ID:
		res += "IonTechnology"
	case PlasmaTechnology.ID:
		res += "PlasmaTechnology"
	case IntergalacticResearchNetwork.ID:
		res += "IntergalacticResearchNetwork"
	case Astrophysics.ID:
		res += "Astrophysics"
	case GravitonTechnology.ID:
		res += "GravitonTechnology"
	default:
		res += "Invalid" + "(" + strconv.Itoa(int(o)) + ")"
	}
	return res
}

// IsFacility ...
func (o ID) IsFacility() bool {
	return o == AllianceDepot.ID ||
		o == RoboticsFactory.ID ||
		o == Shipyard.ID ||
		o == ResearchLab.ID ||
		o == MissileSilo.ID ||
		o == NaniteFactory.ID ||
		o == Terraformer.ID ||
		o == SpaceDock.ID
}

// IsResourceBuilding ...
func (o ID) IsResourceBuilding() bool {
	return o == MetalMine.ID ||
		o == CrystalMine.ID ||
		o == DeuteriumSynthesizer.ID ||
		o == SolarPlant.ID ||
		o == FusionReactor.ID ||
		o == MetalStorage.ID ||
		o == CrystalStorage.ID ||
		o == DeuteriumTank.ID ||
		o == ShieldedMetalDen.ID ||
		o == UndergroundCrystalDen.ID ||
		o == SeabedDeuteriumDen.ID
}

// IsBuilding ...
func (o ID) IsBuilding() bool {
	return o.IsResourceBuilding() || o.IsFacility()
}

// IsTech ...
func (o ID) IsTech() bool {
	return o == EspionageTechnology.ID ||
		o == ComputerTechnology.ID ||
		o == WeaponsTechnology.ID ||
		o == ShieldingTechnology.ID ||
		o == ArmourTechnology.ID ||
		o == EnergyTechnology.ID ||
		o == HyperspaceTechnology.ID ||
		o == CombustionDrive.ID ||
		o == ImpulseDrive.ID ||
		o == HyperspaceDrive.ID ||
		o == LaserTechnology.ID ||
		o == IonTechnology.ID ||
		o == PlasmaTechnology.ID ||
		o == IntergalacticResearchNetwork.ID ||
		o == Astrophysics.ID ||
		o == GravitonTechnology.ID
}

// IsDefense ...
func (o ID) IsDefense() bool {
	return o == RocketLauncher.ID ||
		o == LightLaser.ID ||
		o == HeavyLaser.ID ||
		o == GaussCannon.ID ||
		o == IonCannon.ID ||
		o == PlasmaTurret.ID ||
		o == SmallShieldDome.ID ||
		o == LargeShieldDome.ID ||
		o == AntiBallisticMissiles.ID ||
		o == InterplanetaryMissiles.ID
}

// IsShip ...
func (o ID) IsShip() bool {
	return o == SmallCargo.ID ||
		o == LargeCargo.ID ||
		o == LightFighter.ID ||
		o == HeavyFighter.ID ||
		o == Cruiser.ID ||
		o == Battleship.ID ||
		o == ColonyShip.ID ||
		o == Recycler.ID ||
		o == EspionageProbe.ID ||
		o == Bomber.ID ||
		o == SolarSatellite.ID ||
		o == Destroyer.ID ||
		o == Deathstar.ID ||
		o == Battlecruiser.ID
}
