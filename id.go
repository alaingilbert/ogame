package ogame

import "strconv"

// ID ...
type ID int

// String ...
func (o ID) String() string {
	res := ""
	switch o {
	case AllianceDepot:
		res += "AllianceDepot"
	case RoboticsFactory:
		res += "RoboticsFactory"
	case Shipyard:
		res += "Shipyard"
	case ResearchLab:
		res += "ResearchLab"
	case MissileSilo:
		res += "MissileSilo"
	case NaniteFactory:
		res += "NaniteFactory"
	case Terraformer:
		res += "Terraformer"
	case SpaceDock:
		res += "SpaceDock"
	case MetalMine:
		res += "MetalMine"
	case CrystalMine:
		res += "CrystalMine"
	case DeuteriumSynthesizer:
		res += "DeuteriumSynthesizer"
	case SolarPlant:
		res += "SolarPlant"
	case FusionReactor:
		res += "FusionReactor"
	case MetalStorage:
		res += "MetalStorage"
	case CrystalStorage:
		res += "CrystalStorage"
	case DeuteriumTank:
		res += "DeuteriumTank"
	case ShieldedMetalDen:
		res += "ShieldedMetalDen"
	case UndergroundCrystalDen:
		res += "UndergroundCrystalDen"
	case SeabedDeuteriumDen:
		res += "SeabedDeuteriumDen"
	case RocketLauncher:
		res += "RocketLauncher"
	case LightLaser:
		res += "LightLaser"
	case HeavyLaser:
		res += "HeavyLaser"
	case GaussCannon:
		res += "GaussCannon"
	case IonCannon:
		res += "IonCannon"
	case PlasmaTurret:
		res += "PlasmaTurret"
	case SmallShieldDome:
		res += "SmallShieldDome"
	case LargeShieldDome:
		res += "LargeShieldDome"
	case AntiBallisticMissiles:
		res += "AntiBallisticMissiles"
	case InterplanetaryMissiles:
		res += "InterplanetaryMissiles"
	case SmallCargo:
		res += "SmallCargo"
	case LargeCargo:
		res += "LargeCargo"
	case LightFighter:
		res += "LightFighter"
	case HeavyFighter:
		res += "HeavyFighter"
	case Cruiser:
		res += "Cruiser"
	case Battleship:
		res += "Battleship"
	case ColonyShip:
		res += "ColonyShip"
	case Recycler:
		res += "Recycler"
	case EspionageProbe:
		res += "EspionageProbe"
	case Bomber:
		res += "Bomber"
	case SolarSatellite:
		res += "SolarSatellite"
	case Destroyer:
		res += "Destroyer"
	case Deathstar:
		res += "Deathstar"
	case Battlecruiser:
		res += "Battlecruiser"
	case EspionageTechnology:
		res += "EspionageTechnology"
	case ComputerTechnology:
		res += "ComputerTechnology"
	case WeaponsTechnology:
		res += "WeaponsTechnology"
	case ShieldingTechnology:
		res += "ShieldingTechnology"
	case ArmourTechnology:
		res += "ArmourTechnology"
	case EnergyTechnology:
		res += "EnergyTechnology"
	case HyperspaceTechnology:
		res += "HyperspaceTechnology"
	case CombustionDrive:
		res += "CombustionDrive"
	case ImpulseDrive:
		res += "ImpulseDrive"
	case HyperspaceDrive:
		res += "HyperspaceDrive"
	case LaserTechnology:
		res += "LaserTechnology"
	case IonTechnology:
		res += "IonTechnology"
	case PlasmaTechnology:
		res += "PlasmaTechnology"
	case IntergalacticResearchNetwork:
		res += "IntergalacticResearchNetwork"
	case Astrophysics:
		res += "Astrophysics"
	case GravitonTechnology:
		res += "GravitonTechnology"
	default:
		res += "Invalid" + "(" + strconv.Itoa(int(o)) + ")"
	}
	return res
}

// IsFacility ...
func (o ID) IsFacility() bool {
	return o == AllianceDepot ||
		o == RoboticsFactory ||
		o == Shipyard ||
		o == ResearchLab ||
		o == MissileSilo ||
		o == NaniteFactory ||
		o == Terraformer ||
		o == SpaceDock
}

// IsResourceBuilding ...
func (o ID) IsResourceBuilding() bool {
	return o == MetalMine ||
		o == CrystalMine ||
		o == DeuteriumSynthesizer ||
		o == SolarPlant ||
		o == FusionReactor ||
		o == MetalStorage ||
		o == CrystalStorage ||
		o == DeuteriumTank ||
		o == ShieldedMetalDen ||
		o == UndergroundCrystalDen ||
		o == SeabedDeuteriumDen
}

// IsBuilding ...
func (o ID) IsBuilding() bool {
	return o.IsResourceBuilding() || o.IsFacility()
}

// IsTech ...
func (o ID) IsTech() bool {
	return o == EspionageTechnology ||
		o == ComputerTechnology ||
		o == WeaponsTechnology ||
		o == ShieldingTechnology ||
		o == ArmourTechnology ||
		o == EnergyTechnology ||
		o == HyperspaceTechnology ||
		o == CombustionDrive ||
		o == ImpulseDrive ||
		o == HyperspaceDrive ||
		o == LaserTechnology ||
		o == IonTechnology ||
		o == PlasmaTechnology ||
		o == IntergalacticResearchNetwork ||
		o == Astrophysics ||
		o == GravitonTechnology
}

// IsDefense ...
func (o ID) IsDefense() bool {
	return o == RocketLauncher ||
		o == LightLaser ||
		o == HeavyLaser ||
		o == GaussCannon ||
		o == IonCannon ||
		o == PlasmaTurret ||
		o == SmallShieldDome ||
		o == LargeShieldDome ||
		o == AntiBallisticMissiles ||
		o == InterplanetaryMissiles
}

// IsShip ...
func (o ID) IsShip() bool {
	return o == SmallCargo ||
		o == LargeCargo ||
		o == LightFighter ||
		o == HeavyFighter ||
		o == Cruiser ||
		o == Battleship ||
		o == ColonyShip ||
		o == Recycler ||
		o == EspionageProbe ||
		o == Bomber ||
		o == SolarSatellite ||
		o == Destroyer ||
		o == Deathstar ||
		o == Battlecruiser
}
