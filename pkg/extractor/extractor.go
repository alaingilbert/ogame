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
	ExtractLifeformTypeFromDoc(*goquery.Document) ogame.LifeformType
	ExtractAdmiralFromDoc(*goquery.Document) bool
	ExtractBodyIDFromDoc(*goquery.Document) string
	ExtractCelestialFromDoc(*goquery.Document, any) (ogame.Celestial, error)
	ExtractCelestialsFromDoc(*goquery.Document) ([]ogame.Celestial, error)
	ExtractCharacterClassFromDoc(*goquery.Document) (ogame.CharacterClass, error)
	ExtractCommanderFromDoc(*goquery.Document) bool
	ExtractEngineerFromDoc(*goquery.Document) bool
	ExtractGeologistFromDoc(*goquery.Document) bool
	ExtractIsInVacationFromDoc(*goquery.Document) bool
	ExtractIsMobileFromDoc(*goquery.Document) bool
	ExtractMoonFromDoc(*goquery.Document, any) (ogame.Moon, error)
	ExtractMoonsFromDoc(*goquery.Document) []ogame.Moon
	ExtractOGameSessionFromDoc(*goquery.Document) string
	ExtractOgameTimestampFromDoc(*goquery.Document) int64
	ExtractPlanetFromDoc(*goquery.Document, any) (ogame.Planet, error)
	ExtractPlanetIDFromDoc(*goquery.Document) (ogame.CelestialID, error)
	ExtractPlanetTypeFromDoc(*goquery.Document) (ogame.CelestialType, error)
	ExtractPlanetsFromDoc(*goquery.Document) []ogame.Planet
	ExtractResourcesDetailsFromFullPageFromDoc(*goquery.Document) ogame.ResourcesDetails
	ExtractResourcesFromDoc(*goquery.Document) ogame.Resources
	ExtractServerTimeFromDoc(*goquery.Document) (time.Time, error)
	ExtractTechnocratFromDoc(*goquery.Document) bool
	ExtractColoniesFromDoc(*goquery.Document) (int64, int64)
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
	ExtractConstructions(pageHTML []byte) (ogame.Constructions, error)
	ExtractDMCosts(pageHTML []byte) (ogame.DMCosts, error)
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractUserInfos(pageHTML []byte) (ogame.UserInfos, error)
}

type OverviewExtractorDoc interface {
	ExtractOverviewProductionFromDoc(*goquery.Document) ([]ogame.Quantifiable, error)
	ExtractCharacterClassFromDoc(*goquery.Document) (ogame.CharacterClass, error)
}

type OverviewExtractorBytesDoc interface {
	OverviewExtractorBytes
	OverviewExtractorDoc
}

type FleetsExtractorBytes interface {
	ExtractSlots(pageHTML []byte) (ogame.Slots, error)
}

type FleetsExtractorDoc interface {
	ExtractSlotsFromDoc(*goquery.Document) (ogame.Slots, error)
}

type MovementExtractorBytes interface {
	FleetsExtractorBytes
	ExtractFleets(pageHTML []byte) ([]ogame.Fleet, error)
}

type MovementExtractorDoc interface {
	FleetsExtractorDoc
	ExtractFleetsFromDoc(*goquery.Document) ([]ogame.Fleet, error)
}

type MovementExtractorBytesDoc interface {
	MovementExtractorBytes
	MovementExtractorDoc
}

type FleetDispatchExtractorBytes interface {
	FleetsExtractorBytes
	ExtractFleet1Ships(pageHTML []byte) (ogame.ShipsInfos, error)
}

type FleetDispatchExtractorDoc interface {
	FleetsExtractorBytes
	ExtractFleet1ShipsFromDoc(*goquery.Document) (ogame.ShipsInfos, error)
	ExtractFleetDispatchACSFromDoc(*goquery.Document) []ogame.ACSValues
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
	ExtractProductionFromDoc(*goquery.Document) ([]ogame.Quantifiable, error)
	ExtractShipsFromDoc(*goquery.Document) (ogame.ShipsInfos, error)
}

type ShipyardExtractorBytesDoc interface {
	ShipyardExtractorBytes
	ShipyardExtractorDoc
}

type ResearchExtractorBytes interface {
	ExtractResearch(pageHTML []byte) (ogame.Researches, error)
	ExtractUpgradeToken(pageHTML []byte) (string, error)
}

type ResearchExtractorDoc interface {
	ExtractResearchFromDoc(*goquery.Document) ogame.Researches
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
	ExtractFacilitiesFromDoc(*goquery.Document) (ogame.Facilities, error)
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
	ExtractPopopsCombatreportFromDoc(*goquery.Document) bool
	ExtractPreferences(pageHTML []byte) (ogame.Preferences, error)
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) (bool, error)
	ExtractSpioAnz(pageHTML []byte) (int64, error)
}

type PreferencesExtractorDoc interface {
	ExtractActivateAutofocusFromDoc(*goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(*goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(*goquery.Document) bool
	ExtractAuctioneerNotificationsFromDoc(*goquery.Document) bool
	ExtractDisableChatBarFromDoc(*goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(*goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(*goquery.Document) bool
	ExtractEventsShowFromDoc(*goquery.Document) int64
	ExtractMobileVersionFromDoc(*goquery.Document) bool
	ExtractMsgResultsPerPageFromDoc(*goquery.Document) int64
	ExtractNotifAccountFromDoc(*goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(*goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(*goquery.Document) bool
	ExtractNotifAuctionsFromDoc(*goquery.Document) bool
	ExtractNotifBuildListFromDoc(*goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(*goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(*goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(*goquery.Document) bool
	ExtractPopupsNoticesFromDoc(*goquery.Document) bool
	ExtractPreferencesFromDoc(*goquery.Document) ogame.Preferences
	ExtractPreserveSystemOnPlanetChangeFromDoc(*goquery.Document) bool
	ExtractShowActivityMinutesFromDoc(*goquery.Document) bool
	ExtractShowDetailOverlayFromDoc(*goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(*goquery.Document) bool
	ExtractSortOrderFromDoc(*goquery.Document) int64
	ExtractSortSettingFromDoc(*goquery.Document) int64
	ExtractSpioAnzFromDoc(*goquery.Document) int64
	ExtractSpioReportPicturesFromDoc(*goquery.Document) bool
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
	ExtractDefenseFromDoc(*goquery.Document) (ogame.DefensesInfos, error)
}

type DefensesExtractorBytesDoc interface {
	DefensesExtractorBytes
	DefensesExtractorDoc
}

type EventListExtractorBytes interface {
	ExtractAttacks(pageHTML []byte, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractFleetsFromEventList(pageHTML []byte) ([]ogame.Fleet, error)
}

type EventListExtractorDoc interface {
	ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractFleetsFromEventListFromDoc(*goquery.Document) []ogame.Fleet
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
	ExtractCombatReportMessagesFromDoc(*goquery.Document) ([]ogame.CombatReportSummary, int64, error)
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
	ExtractEspionageReportFromDoc(*goquery.Document) (ogame.EspionageReport, error)
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
	ExtractEspionageReportMessageIDsFromDoc(*goquery.Document) ([]ogame.EspionageReportSummary, int64, error)
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
	ExtractExpeditionMessagesFromDoc(*goquery.Document) ([]ogame.ExpeditionMessage, int64, error)
}

type MessagesExpeditionExtractorBytesDoc interface {
	MessagesExpeditionExtractorBytes
	MessagesExpeditionExtractorDoc
}

// FederationExtractorBytes popup when we click to create a union for our attacking fleet
type FederationExtractorBytes interface {
	ExtractFederation(pageHTML []byte) (url.Values, error)
}

// GalaxyPageExtractorBytes galaxy page
type GalaxyPageExtractorBytes interface {
	ExtractAvailableDiscoveries(pageHTML []byte) (int64, error)
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
	ExtractOfferOfTheDayFromDoc(*goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error)
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
	ExtractResourceSettingsFromDoc(*goquery.Document) (ogame.ResourceSettings, string, error)
	ExtractResourcesProductionsFromDoc(*goquery.Document) (ogame.Resources, error)
}

type ResourcesSettingsExtractorBytesDoc interface {
	ResourcesSettingsExtractorBytes
	ResourcesSettingsExtractorDoc
}

type HighscoreExtractorBytes interface {
	ExtractHighscore(pageHTML []byte) (ogame.Highscore, error)
}

type HighscoreExtractorDoc interface {
	ExtractHighscoreFromDoc(*goquery.Document) (ogame.Highscore, error)
}

type HighscoreExtractorBytesDoc interface {
	HighscoreExtractorBytes
	HighscoreExtractorDoc
}

type MissileAttackLayerExtractorBytes interface {
	ExtractIPM(pageHTML []byte) (duration, max int64, token string, err error)
}

type MissileAttackLayerExtractorDoc interface {
	ExtractIPMFromDoc(*goquery.Document) (duration, max int64, token string, err error)
}

type MissileAttackLayerExtractorBytesDoc interface {
	MissileAttackLayerExtractorBytes
	MissileAttackLayerExtractorDoc
}

type JumpGateLayerExtractorBytes interface {
	ExtractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64, error)
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
	ExtractLfBuildingsFromDoc(*goquery.Document) (ogame.LfBuildings, error)
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
	ExtractLfResearchFromDoc(*goquery.Document) (ogame.LfResearches, error)
	ExtractLfSlotsFromDoc(*goquery.Document) [18]ogame.LfSlot
	ExtractArtefactsFromDoc(*goquery.Document) (int64, int64)
}

type LfResearchExtractorBytesDoc interface {
	LfResearchExtractorBytes
	LfResearchExtractorDoc
}

type LfBonusesExtractorBytes interface {
	ExtractLfBonuses(pageHTML []byte) (ogame.LfBonuses, error)
}

type LfBonusesExtractorDoc interface {
	ExtractLfBonusesFromDoc(*goquery.Document) (ogame.LfBonuses, error)
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
	ExtractResourcesBuildingsFromDoc(*goquery.Document) (ogame.ResourcesBuildings, error)
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
	ExtractAbandonInformation(*goquery.Document) (abandonToken string, token string)
}

type TechnologyDetailsExtractorBytes interface {
	ExtractTechnologyDetails(pageHTML []byte) (ogame.TechnologyDetails, error)
	ExtractTearDownButtonEnabled(pageHTML []byte) (bool, error)
}

type TechnologyDetailsExtractorDoc interface {
	ExtractTearDownButtonEnabledFromDoc(*goquery.Document) bool
	ExtractTechnologyDetailsFromDoc(*goquery.Document) (ogame.TechnologyDetails, error)
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
	SetLocation(*time.Location)
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

	ExtractCoord(v string) ogame.Coordinate
	ExtractHiddenFields(pageHTML []byte) (url.Values, error)

	ExtractHiddenFieldsFromDoc(*goquery.Document) url.Values
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
