package extractor

import (
	v10 "github.com/alaingilbert/ogame/pkg/extractor/v10"
	v104 "github.com/alaingilbert/ogame/pkg/extractor/v104"
	v11 "github.com/alaingilbert/ogame/pkg/extractor/v11"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_13_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_15_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_9_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v12_0_0"
	v71 "github.com/alaingilbert/ogame/pkg/extractor/v71"
	v8 "github.com/alaingilbert/ogame/pkg/extractor/v8"
	v874 "github.com/alaingilbert/ogame/pkg/extractor/v874"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	v7 "github.com/alaingilbert/ogame/pkg/extractor/v7"
	v9 "github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

type FullPageExtractorBytes interface {
	ExtractAdmiral(pageHTML []byte) (bool, error)
	ExtractAjaxChatToken(pageHTML []byte) (string, error)
	ExtractToken(pageHTML []byte) (string, error)
	ExtractCelestial(pageHTML []byte, v any) (ogame.Celestial, error)
	ExtractCelestials(pageHTML []byte) ([]ogame.Celestial, error)
	ExtractCharacterClass(pageHTML []byte) (ogame.CharacterClass, error)
	ExtractCommander(pageHTML []byte) (bool, error)
	ExtractEngineer(pageHTML []byte) (bool, error)
	ExtractGeologist(pageHTML []byte) (bool, error)
	ExtractIsInVacation(pageHTML []byte) (bool, error)
	ExtractIsMobile(pageHTML []byte) (bool, error)
	ExtractLifeformEnabled(pageHTML []byte) bool
	ExtractMoon(pageHTML []byte, v any) (ogame.Moon, error)
	ExtractMoons(pageHTML []byte) ([]ogame.Moon, error)
	ExtractOGameTimestampFromBytes(pageHTML []byte) (int64, error)
	ExtractOgameTimestamp(pageHTML []byte) (int64, error)
	ExtractPlanet(pageHTML []byte, v any) (ogame.Planet, error)
	ExtractPlanetCoordinate(pageHTML []byte) (ogame.Coordinate, error)
	ExtractPlanetID(pageHTML []byte) (ogame.CelestialID, error)
	ExtractPlanetType(pageHTML []byte) (ogame.CelestialType, error)
	ExtractPlanets(pageHTML []byte) ([]ogame.Planet, error)
	ExtractResources(pageHTML []byte) (ogame.Resources, error)
	ExtractResourcesDetailsFromFullPage(pageHTML []byte) (ogame.ResourcesDetails, error)
	ExtractServerTime(pageHTML []byte) (time.Time, error)
	ExtractTechnocrat(pageHTML []byte) (bool, error)
}

type FullPageExtractorDoc interface {
	ExtractLifeformTypeFromDoc(doc *goquery.Document) ogame.LifeformType
	ExtractAdmiralFromDoc(doc *goquery.Document) bool
	ExtractBodyIDFromDoc(doc *goquery.Document) string
	ExtractCelestialFromDoc(doc *goquery.Document, v any) (ogame.Celestial, error)
	ExtractCelestialsFromDoc(doc *goquery.Document) ([]ogame.Celestial, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (ogame.CharacterClass, error)
	ExtractCommanderFromDoc(doc *goquery.Document) bool
	ExtractEngineerFromDoc(doc *goquery.Document) bool
	ExtractGeologistFromDoc(doc *goquery.Document) bool
	ExtractIsInVacationFromDoc(doc *goquery.Document) bool
	ExtractIsMobileFromDoc(doc *goquery.Document) bool
	ExtractMoonFromDoc(doc *goquery.Document, v any) (ogame.Moon, error)
	ExtractMoonsFromDoc(doc *goquery.Document) []ogame.Moon
	ExtractOGameSessionFromDoc(doc *goquery.Document) string
	ExtractOgameTimestampFromDoc(doc *goquery.Document) int64
	ExtractPlanetFromDoc(doc *goquery.Document, v any) (ogame.Planet, error)
	ExtractPlanetIDFromDoc(doc *goquery.Document) (ogame.CelestialID, error)
	ExtractPlanetTypeFromDoc(doc *goquery.Document) (ogame.CelestialType, error)
	ExtractPlanetsFromDoc(doc *goquery.Document) []ogame.Planet
	ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails
	ExtractResourcesFromDoc(doc *goquery.Document) ogame.Resources
	ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error)
	ExtractTechnocratFromDoc(doc *goquery.Document) bool
	ExtractColoniesFromDoc(doc *goquery.Document) (int64, int64)
}

type FullPageExtractorBytesDoc interface {
	FullPageExtractorBytes
	FullPageExtractorDoc
}

type OverviewExtractorBytes interface {
	ExtractActiveItems(pageHTML []byte) ([]ogame.ActiveItem, error)
	ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCancelFleetToken(pageHTML []byte, fleetID ogame.FleetID) (string, error)
	ExtractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error)
	ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCharacterClass(pageHTML []byte) (ogame.CharacterClass, error)
	ExtractConstructions(pageHTML []byte) ogame.Constructions
	ExtractDMCosts(pageHTML []byte) (ogame.DMCosts, error)
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractUserInfos(pageHTML []byte) (ogame.UserInfos, error)
}

type OverviewExtractorDoc interface {
	ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (ogame.CharacterClass, error)
}

type OverviewExtractorBytesDoc interface {
	OverviewExtractorBytes
	OverviewExtractorDoc
}

type FleetsExtractorBytes interface {
	ExtractSlots(pageHTML []byte) (ogame.Slots, error)
}

type FleetsExtractorDoc interface {
	ExtractSlotsFromDoc(doc *goquery.Document) (ogame.Slots, error)
}

type MovementExtractorBytes interface {
	FleetsExtractorBytes
	ExtractFleets(pageHTML []byte) (res []ogame.Fleet)
}

type MovementExtractorDoc interface {
	FleetsExtractorDoc
	ExtractFleetsFromDoc(doc *goquery.Document) (res []ogame.Fleet)
}

type MovementExtractorBytesDoc interface {
	MovementExtractorBytes
	MovementExtractorDoc
}

type FleetDispatchExtractorBytes interface {
	FleetsExtractorBytes
	ExtractFleet1Ships(pageHTML []byte) ogame.ShipsInfos
}

type FleetDispatchExtractorDoc interface {
	FleetsExtractorBytes
	ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ogame.ShipsInfos)
	ExtractFleetDispatchACSFromDoc(doc *goquery.Document) []ogame.ACSValues
}

type FleetDispatchExtractorBytesDoc interface {
	FleetDispatchExtractorBytes
	FleetDispatchExtractorDoc
}

type ShipyardExtractorBytes interface {
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error)
	ExtractShips(pageHTML []byte) (ogame.ShipsInfos, error)
	ExtractUpgradeToken(pageHTML []byte) (string, error)
}

type ShipyardExtractorDoc interface {
	ExtractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error)
	ExtractShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error)
}

type ShipyardExtractorBytesDoc interface {
	ShipyardExtractorBytes
	ShipyardExtractorDoc
}

type ResearchExtractorBytes interface {
	ExtractResearch(pageHTML []byte) ogame.Researches
	ExtractUpgradeToken(pageHTML []byte) (string, error)
}

type ResearchExtractorDoc interface {
	ExtractResearchFromDoc(doc *goquery.Document) ogame.Researches
}

type ResearchExtractorBytesDoc interface {
	ResearchExtractorBytes
	ResearchExtractorDoc
}

type FacilitiesExtractorBytes interface {
	ExtractFacilities(pageHTML []byte) (ogame.Facilities, error)
	ExtractTearDownToken(pageHTML []byte) (string, error)
}

type FacilitiesExtractorDoc interface {
	ExtractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error)
}

type FacilitiesExtractorBytesDoc interface {
	FacilitiesExtractorBytes
	FacilitiesExtractorDoc
}

type PhalanxExtractorBytes interface {
	ExtractPhalanx(pageHTML []byte) ([]ogame.PhalanxFleet, error)
	ExtractPhalanxNewToken(pageHTML []byte) (string, error)
}

type PreferencesExtractorBytes interface {
	ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool
	ExtractPreferences(pageHTML []byte) ogame.Preferences
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool
	ExtractSpioAnz(pageHTML []byte) int64
}

type PreferencesExtractorDoc interface {
	ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool
	ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool
	ExtractDisableChatBarFromDoc(doc *goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool
	ExtractEventsShowFromDoc(doc *goquery.Document) int64
	ExtractMobileVersionFromDoc(doc *goquery.Document) bool
	ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64
	ExtractNotifAccountFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool
	ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool
	ExtractNotifBuildListFromDoc(doc *goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool
	ExtractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences
	ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool
	ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool
	ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool
	ExtractSortOrderFromDoc(doc *goquery.Document) int64
	ExtractSortSettingFromDoc(doc *goquery.Document) int64
	ExtractSpioAnzFromDoc(doc *goquery.Document) int64
	ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool
}

type PreferencesExtractorBytesDoc interface {
	PreferencesExtractorBytes
	PreferencesExtractorDoc
}

type DefensesExtractorBytes interface {
	ExtractDefense(pageHTML []byte) (ogame.DefensesInfos, error)
	ExtractUpgradeToken(pageHTML []byte) (string, error)
}

type DefensesExtractorDoc interface {
	ExtractDefenseFromDoc(doc *goquery.Document) (ogame.DefensesInfos, error)
}

type DefensesExtractorBytesDoc interface {
	DefensesExtractorBytes
	DefensesExtractorDoc
}

type EventListExtractorBytes interface {
	ExtractAttacks(pageHTML []byte, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractFleetsFromEventList(pageHTML []byte) []ogame.Fleet
}

type EventListExtractorDoc interface {
	ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []ogame.Fleet
}

type EventListExtractorBytesDoc interface {
	EventListExtractorBytes
	EventListExtractorDoc
}

type TraderAuctioneerExtractorBytes interface {
	ExtractAllResources(pageHTML []byte) (map[ogame.CelestialID]ogame.Resources, error)
	ExtractAuction(pageHTML []byte) (ogame.Auction, error)
}

type AllianceOverviewExtractorBytes interface {
	ExtractAllianceClass(pageHTML []byte) (ogame.AllianceClass, error)
}

// BuffActivationExtractorBytes BuffActivation is the popups that shows up when clicking the icon
// to activate an item on the overview page.
type BuffActivationExtractorBytes interface {
	ExtractBuffActivation(pageHTML []byte) (string, []ogame.Item, error)
}

type MessagesCombatReportExtractorBytes interface {
	ExtractCombatReportMessagesSummary(pageHTML []byte) ([]ogame.CombatReportSummary, int64, error)
}

type MessagesCombatReportExtractorDoc interface {
	ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error)
}

type MessagesCombatReportExtractorBytesDoc interface {
	MessagesCombatReportExtractorBytes
	MessagesCombatReportExtractorDoc
}

// DestroyRocketsExtractorBytes popups that shows up when clicking to destroy rockets on the defenses page.
type DestroyRocketsExtractorBytes interface {
	ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error)
}

type EmpireExtractorBytes interface {
	ExtractEmpire(pageHTML []byte) ([]ogame.EmpireCelestial, error)
	ExtractEmpireJSON(pageHTML []byte) (any, error)
}

// EspionageReportExtractorBytes popup that shows the full espionage report
type EspionageReportExtractorBytes interface {
	ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error)
}

type EspionageReportExtractorDoc interface {
	ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error)
}

type EspionageReportExtractorBytesDoc interface {
	EspionageReportExtractorBytes
	EspionageReportExtractorDoc
}

// MessagesEspionageReportExtractorBytes ajax page that display all espionage reports summaries
type MessagesEspionageReportExtractorBytes interface {
	ExtractEspionageReportMessageIDs(pageHTML []byte) ([]ogame.EspionageReportSummary, int64, error)
}

type MessagesEspionageReportExtractorDoc interface {
	ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64, error)
}

type MessagesEspionageReportExtractorBytesDoc interface {
	MessagesEspionageReportExtractorBytes
	MessagesEspionageReportExtractorDoc
}

// MessagesExpeditionExtractorBytes ajax page that display all expedition messages
type MessagesExpeditionExtractorBytes interface {
	ExtractExpeditionMessages(pageHTML []byte) ([]ogame.ExpeditionMessage, int64, error)
}

type MessagesExpeditionExtractorDoc interface {
	ExtractExpeditionMessagesFromDoc(doc *goquery.Document) ([]ogame.ExpeditionMessage, int64, error)
}

type MessagesExpeditionExtractorBytesDoc interface {
	MessagesExpeditionExtractorBytes
	MessagesExpeditionExtractorDoc
}

// FederationExtractorBytes popup when we click to create a union for our attacking fleet
type FederationExtractorBytes interface {
	ExtractFederation(pageHTML []byte) url.Values
}

// GalaxyPageExtractorBytes galaxy page
type GalaxyPageExtractorBytes interface {
	ExtractAvailableDiscoveries(pageHTML []byte) int64
}

// GalaxyExtractorBytes ajax page containing galaxy information in galaxy page
type GalaxyExtractorBytes interface {
	ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (ogame.SystemInfos, error)
}

// FetchResourcesExtractorBytes "fetchResources" ajax page
type FetchResourcesExtractorBytes interface {
	ExtractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error)
}

// TraderImportExportExtractorBytes ajax page Merchant -> Import/Export
type TraderImportExportExtractorBytes interface {
	ExtractOfferOfTheDay(pageHTML []byte) (int64, string, ogame.PlanetResources, ogame.Multiplier, error)
}

type TraderImportExportExtractorDoc interface {
	ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error)
}

// FetchTechsExtractorBytes ajax page fetchTechs
type FetchTechsExtractorBytes interface {
	ExtractTechs(pageHTML []byte) (ogame.Techs, error)
}

type ResourcesSettingsExtractorBytes interface {
	ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, string, error)
	ExtractResourcesProductions(pageHTML []byte) (ogame.Resources, error)
}

type ResourcesSettingsExtractorDoc interface {
	ExtractResourceSettingsFromDoc(doc *goquery.Document) (ogame.ResourceSettings, string, error)
	ExtractResourcesProductionsFromDoc(doc *goquery.Document) (ogame.Resources, error)
}

type ResourcesSettingsExtractorBytesDoc interface {
	ResourcesSettingsExtractorBytes
	ResourcesSettingsExtractorDoc
}

type HighscoreExtractorBytes interface {
	ExtractHighscore(pageHTML []byte) (ogame.Highscore, error)
}

type HighscoreExtractorDoc interface {
	ExtractHighscoreFromDoc(doc *goquery.Document) (ogame.Highscore, error)
}

type HighscoreExtractorBytesDoc interface {
	HighscoreExtractorBytes
	HighscoreExtractorDoc
}

type MissileAttackLayerExtractorBytes interface {
	ExtractIPM(pageHTML []byte) (duration, max int64, token string)
}

type MissileAttackLayerExtractorDoc interface {
	ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string)
}

type MissileAttackLayerExtractorBytesDoc interface {
	MissileAttackLayerExtractorBytes
	MissileAttackLayerExtractorDoc
}

type JumpGateLayerExtractorBytes interface {
	ExtractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64)
}

// MessagesMarketplaceExtractorBytes marketplace was removed from the game
type MessagesMarketplaceExtractorBytes interface {
	ExtractMarketplaceMessages(pageHTML []byte) ([]ogame.MarketplaceMessage, int64, error)
}

type LfBuildingsExtractorBytes interface {
	ExtractUpgradeToken(pageHTML []byte) (string, error)
	ExtractLfBuildings(pageHTML []byte) (ogame.LfBuildings, error)
}

type LfBuildingsExtractorDoc interface {
	ExtractLfBuildingsFromDoc(doc *goquery.Document) (ogame.LfBuildings, error)
}

type LfBuildingsExtractorBytesDoc interface {
	LfBuildingsExtractorBytes
	LfBuildingsExtractorDoc
}

type LfResearchExtractorBytes interface {
	ExtractUpgradeToken(pageHTML []byte) (string, error)
	ExtractLfResearch(pageHTML []byte) (ogame.LfResearches, error)
}

type LfResearchExtractorDoc interface {
	ExtractLfResearchFromDoc(doc *goquery.Document) (ogame.LfResearches, error)
	ExtractLfSlotsFromDoc(doc *goquery.Document) [18]ogame.LfSlot
	ExtractArtefactsFromDoc(doc *goquery.Document) (int64, int64)
}

type LfResearchExtractorBytesDoc interface {
	LfResearchExtractorBytes
	LfResearchExtractorDoc
}

type LfBonusesExtractorBytes interface {
	ExtractLfBonuses(pageHTML []byte) (ogame.LfBonuses, error)
}

type LfBonusesExtractorDoc interface {
	ExtractLfBonusesFromDoc(doc *goquery.Document) (ogame.LfBonuses, error)
}

type LfBonusesExtractorBytesDoc interface {
	LfBonusesExtractorBytes
	LfBonusesExtractorDoc
}

// ResourcesBuildingsExtractorBytes supplies page
type ResourcesBuildingsExtractorBytes interface {
	ExtractResourcesBuildings(pageHTML []byte) (ogame.ResourcesBuildings, error)
	ExtractTearDownToken(pageHTML []byte) (string, error)
	ExtractUpgradeToken(pageHTML []byte) (string, error)
}

type ResourcesBuildingsExtractorDoc interface {
	ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ogame.ResourcesBuildings, error)
}

type ResourcesBuildingsExtractorBytesDoc interface {
	ResourcesBuildingsExtractorBytes
	ResourcesBuildingsExtractorDoc
}

// PremiumExtractorBytes ajax page when click to buy an officer
type PremiumExtractorBytes interface {
	ExtractPremiumToken(pageHTML []byte, days int64) (token string, err error)
}

type PlanetLayerExtractorDoc interface {
	ExtractAbandonInformation(doc *goquery.Document) (abandonToken string, token string)
}

type TechnologyDetailsExtractorBytes interface {
	ExtractTechnologyDetails(pageHTML []byte) (ogame.TechnologyDetails, error)
	ExtractTearDownButtonEnabled(pageHTML []byte) bool
}

type TechnologyDetailsExtractorDoc interface {
	ExtractTearDownButtonEnabledFromDoc(doc *goquery.Document) bool
	ExtractTechnologyDetailsFromDoc(doc *goquery.Document) (ogame.TechnologyDetails, error)
}

type TechnologyDetailsExtractorBytesDoc interface {
	TechnologyDetailsExtractorBytes
	TechnologyDetailsExtractorDoc
}

// Extractor ...
type Extractor interface {
	GetLanguage() string
	SetLanguage(lang string)
	GetLocation() *time.Location
	SetLocation(loc *time.Location)
	GetLifeformEnabled() bool
	SetLifeformEnabled(lifeformEnabled bool)

	DefensesExtractorBytesDoc
	EspionageReportExtractorBytesDoc
	EventListExtractorBytesDoc
	FacilitiesExtractorBytesDoc
	FleetDispatchExtractorBytesDoc
	FullPageExtractorBytesDoc
	HighscoreExtractorBytesDoc
	LfBuildingsExtractorBytesDoc
	LfResearchExtractorBytesDoc
	LfBonusesExtractorBytesDoc
	MessagesCombatReportExtractorBytesDoc
	MessagesEspionageReportExtractorBytesDoc
	MessagesExpeditionExtractorBytesDoc
	MissileAttackLayerExtractorBytesDoc
	MovementExtractorBytesDoc
	OverviewExtractorBytesDoc
	PreferencesExtractorBytesDoc
	ResearchExtractorBytesDoc
	ResourcesBuildingsExtractorBytesDoc
	ResourcesSettingsExtractorBytesDoc
	ShipyardExtractorBytesDoc
	TechnologyDetailsExtractorBytesDoc

	BuffActivationExtractorBytes
	DestroyRocketsExtractorBytes
	EmpireExtractorBytes
	FederationExtractorBytes
	FetchResourcesExtractorBytes
	FetchTechsExtractorBytes
	GalaxyExtractorBytes
	GalaxyPageExtractorBytes
	JumpGateLayerExtractorBytes
	MessagesMarketplaceExtractorBytes
	PhalanxExtractorBytes
	PremiumExtractorBytes
	TraderAuctioneerExtractorBytes
	TraderImportExportExtractorBytes
	AllianceOverviewExtractorBytes

	PlanetLayerExtractorDoc
	TraderImportExportExtractorDoc

	ExtractCoord(v string) (coord ogame.Coordinate)
	ExtractHiddenFields(pageHTML []byte) (fields url.Values)

	ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values
}

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = (*v6.Extractor)(nil)
var _ Extractor = (*v7.Extractor)(nil)
var _ Extractor = (*v71.Extractor)(nil)
var _ Extractor = (*v8.Extractor)(nil)
var _ Extractor = (*v874.Extractor)(nil)
var _ Extractor = (*v9.Extractor)(nil)
var _ Extractor = (*v10.Extractor)(nil)
var _ Extractor = (*v104.Extractor)(nil)
var _ Extractor = (*v11.Extractor)(nil)
var _ Extractor = (*v11_9_0.Extractor)(nil)
var _ Extractor = (*v11_13_0.Extractor)(nil)
var _ Extractor = (*v11_15_0.Extractor)(nil)
var _ Extractor = (*v12_0_0.Extractor)(nil)
