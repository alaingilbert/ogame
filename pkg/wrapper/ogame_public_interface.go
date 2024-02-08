package wrapper

import (
	"crypto/tls"
	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
	"net/http"
	"net/url"
	"time"
)

// Public interface -----------------------------------------------------------

// Enable enables communications with OGame Server
func (b *OGame) Enable() {
	b.enable()
}

// Disable disables communications with OGame Server
func (b *OGame) Disable() {
	b.disable()
}

// IsEnabled returns true if the bot is enabled, otherwise false
func (b *OGame) IsEnabled() bool {
	return b.isEnabled()
}

// IsLoggedIn returns true if the bot is currently logged-in, otherwise false
func (b *OGame) IsLoggedIn() bool {
	return b.isLoggedInAtom.Load()
}

// IsConnected returns true if the bot is currently connected (communication between the bot and OGame is possible), otherwise false
func (b *OGame) IsConnected() bool {
	return b.isConnectedAtom.Load()
}

// GetDevice get the device used by the bot
func (b *OGame) GetDevice() *device.Device {
	return b.device
}

// GetClient get the http client used by the bot
func (b *OGame) GetClient() *httpclient.Client {
	return b.device.GetClient()
}

// SetClient set the http client used by the bot
func (b *OGame) SetClient(client *httpclient.Client) {
	b.device.SetClient(client)
}

// GetLoginClient get the http client used by the bot for login operations
func (b *OGame) GetLoginClient() *httpclient.Client {
	return b.device.GetClient()
}

// GetPublicIP get the public IP used by the bot
func (b *OGame) GetPublicIP() (string, error) {
	return b.getPublicIP()
}

// ValidateAccount validate a gameforge account
func (b *OGame) ValidateAccount(code string) error {
	return b.validateAccount(code)
}

// OnStateChange register a callback that is notified when the bot state changes
func (b *OGame) OnStateChange(clb func(locked bool, actor string)) {
	b.stateChangeCallbacks = append(b.stateChangeCallbacks, clb)
}

// GetState returns the current bot state
func (b *OGame) GetState() (bool, string) {
	return b.lockedAtom.Load(), b.state
}

// IsLocked returns either or not the bot is currently locked
func (b *OGame) IsLocked() bool {
	return b.lockedAtom.Load()
}

// GetSession get ogame session
func (b *OGame) GetSession() string {
	return b.ogameSession
}

// AddAccount add a new account (server) to your list of accounts
func (b *OGame) AddAccount(number int, lang string) (*gameforge.AddAccountRes, error) {
	return b.addAccount(number, lang)
}

// WithPriority ...
func (b *OGame) WithPriority(priority taskRunner.Priority) Prioritizable {
	return b.taskRunnerInst.WithPriority(priority)
}

// Begin start a transaction. Once this function is called, "Done" must be called to release the lock.
func (b *OGame) Begin() Prioritizable {
	return b.WithPriority(taskRunner.Normal).Begin()
}

// BeginNamed begins a new transaction with a name. "Done" must be called to release the lock.
func (b *OGame) BeginNamed(name string) Prioritizable {
	return b.WithPriority(taskRunner.Normal).BeginNamed(name)
}

// SetInitiator ...
func (b *OGame) SetInitiator(initiator string) Prioritizable {
	return nil
}

// Done ...
func (b *OGame) Done() {}

// Tx locks the bot during the transaction and ensure the lock is released afterward
func (b *OGame) Tx(clb func(tx Prioritizable) error) error {
	return b.WithPriority(taskRunner.Normal).Tx(clb)
}

// GetServer get ogame server information that the bot is connected to
func (b *OGame) GetServer() gameforge.Server {
	return b.server
}

// GetServerData get ogame server data information that the bot is connected to
func (b *OGame) GetServerData() gameforge.ServerData {
	return b.serverData
}

// ServerURL get the ogame server specific url
func (b *OGame) ServerURL() string {
	return b.serverURL
}

// GetLanguage get ogame server language
func (b *OGame) GetLanguage() string {
	return b.language
}

// LoginWithBearerToken to ogame server reusing existing token
func (b *OGame) LoginWithBearerToken(token string) (bool, error) {
	return b.WithPriority(taskRunner.Normal).LoginWithBearerToken(token)
}

// LoginWithExistingCookies to ogame server reusing existing cookies
func (b *OGame) LoginWithExistingCookies() (bool, error) {
	return b.WithPriority(taskRunner.Normal).LoginWithExistingCookies()
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *OGame) Login() error {
	return b.WithPriority(taskRunner.Normal).Login()
}

// Logout the bot from ogame server
func (b *OGame) Logout() { b.WithPriority(taskRunner.Normal).Logout() }

// BytesDownloaded returns the amount of bytes downloaded
func (b *OGame) BytesDownloaded() int64 {
	return b.device.GetClient().BytesDownloaded()
}

// BytesUploaded returns the amount of bytes uploaded
func (b *OGame) BytesUploaded() int64 {
	return b.device.GetClient().BytesUploaded()
}

// GetUniverseName get the name of the universe the bot is playing into
func (b *OGame) GetUniverseName() string {
	return b.Universe
}

// GetUsername get the username that was used to login on ogame server
func (b *OGame) GetUsername() string {
	return b.Username
}

// GetResearchSpeed gets the research speed
func (b *OGame) GetResearchSpeed() int64 {
	return b.serverData.ResearchDurationDivisor
}

// GetNbSystems gets the number of systems
func (b *OGame) GetNbSystems() int64 {
	return b.serverData.Systems
}

// GetUniverseSpeed shortcut to get ogame universe speed
func (b *OGame) GetUniverseSpeed() int64 {
	return b.getUniverseSpeed()
}

// GetUniverseSpeedFleet shortcut to get ogame universe speed fleet
func (b *OGame) GetUniverseSpeedFleet() int64 {
	return b.getUniverseSpeedFleet()
}

// IsPioneers either or not the bot use lobby-pioneers
func (b *OGame) IsPioneers() bool {
	return b.lobby == gameforge.LobbyPioneers
}

// IsDonutGalaxy shortcut to get ogame galaxy donut config
func (b *OGame) IsDonutGalaxy() bool {
	return b.isDonutGalaxy()
}

// IsDonutSystem shortcut to get ogame system donut config
func (b *OGame) IsDonutSystem() bool {
	return b.isDonutSystem()
}

// ConstructionTime get duration to build something
func (b *OGame) ConstructionTime(id ogame.ID, nbr int64, facilities ogame.Facilities) time.Duration {
	return b.constructionTime(id, nbr, facilities)
}

// FleetDeutSaveFactor returns the fleet deut save factor
func (b *OGame) FleetDeutSaveFactor() float64 {
	return b.serverData.GlobalDeuteriumSaveFactor
}

// GetPageContent gets the html for a specific ogame page
func (b *OGame) GetPageContent(vals url.Values) ([]byte, error) {
	return b.WithPriority(taskRunner.Normal).GetPageContent(vals)
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *OGame) PostPageContent(vals, payload url.Values) ([]byte, error) {
	return b.WithPriority(taskRunner.Normal).PostPageContent(vals, payload)
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *OGame) IsUnderAttack(opts ...Option) (bool, error) {
	return b.WithPriority(taskRunner.Normal).IsUnderAttack(opts...)
}

// GetCachedPlayer returns cached player infos
func (b *OGame) GetCachedPlayer() ogame.UserInfos {
	return b.Player
}

// GetCachedPreferences returns cached preferences
func (b *OGame) GetCachedPreferences() ogame.Preferences {
	return b.CachedPreferences
}

// SetVacationMode puts account in vacation mode
func (b *OGame) SetVacationMode() error {
	return b.WithPriority(taskRunner.Normal).SetVacationMode()
}

// SetPreferences ...
func (b *OGame) SetPreferences(p ogame.Preferences) error {
	return b.WithPriority(taskRunner.Normal).SetPreferences(p)
}

// SetPreferencesLang ...
func (b *OGame) SetPreferencesLang(lang string) error {
	return b.WithPriority(taskRunner.Normal).SetPreferencesLang(lang)
}

// IsVacationModeEnabled returns either or not the bot is in vacation mode
func (b *OGame) IsVacationModeEnabled() bool {
	return b.isVacationModeEnabled
}

// GetPlanets returns the user planets
func (b *OGame) GetPlanets() ([]Planet, error) {
	return b.WithPriority(taskRunner.Normal).GetPlanets()
}

// GetCachedPlanet return planet from cached value
func (b *OGame) GetCachedPlanet(v IntoPlanet) (Planet, error) {
	return b.getCachedPlanet(v)
}

// GetCachedMoon return moon from cached value
func (b *OGame) GetCachedMoon(v IntoMoon) (Moon, error) {
	return b.getCachedMoon(v)
}

// GetCachedPlanets return planets from cached value
func (b *OGame) GetCachedPlanets() []Planet {
	return b.getCachedPlanets()
}

// GetCachedMoons return moons from cached value
func (b *OGame) GetCachedMoons() []Moon {
	return b.getCachedMoons()
}

// GetCachedCelestials get all cached celestials
func (b *OGame) GetCachedCelestials() []Celestial {
	return b.getCachedCelestials()
}

// GetCachedCelestial return celestial from cached value
func (b *OGame) GetCachedCelestial(v IntoCelestial) (Celestial, error) {
	return b.getCachedCelestial(v)
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *OGame) GetPlanet(v IntoPlanet) (Planet, error) {
	return b.WithPriority(taskRunner.Normal).GetPlanet(v)
}

// GetMoons returns the user moons
func (b *OGame) GetMoons() ([]Moon, error) {
	return b.WithPriority(taskRunner.Normal).GetMoons()
}

// GetMoon gets infos for moonID
func (b *OGame) GetMoon(v IntoMoon) (Moon, error) {
	return b.WithPriority(taskRunner.Normal).GetMoon(v)
}

// GetCelestials get the player's planets & moons
func (b *OGame) GetCelestials() ([]Celestial, error) {
	return b.WithPriority(taskRunner.Normal).GetCelestials()
}

// RecruitOfficer recruit an officer.
// Typ 2: Commander, 3: Admiral, 4: Engineer, 5: Geologist, 6: Technocrat
// Days: 7 or 90
func (b *OGame) RecruitOfficer(typ, days int64) error {
	return b.WithPriority(taskRunner.Normal).RecruitOfficer(typ, days)
}

// Abandon a planet
func (b *OGame) Abandon(v IntoPlanet) error {
	return b.WithPriority(taskRunner.Normal).Abandon(v)
}

// GetCelestial get the player's planet/moon using the coordinate
func (b *OGame) GetCelestial(v IntoCelestial) (Celestial, error) {
	return b.WithPriority(taskRunner.Normal).GetCelestial(v)
}

// ServerVersion returns OGame version
func (b *OGame) ServerVersion() string {
	return b.serverData.Version
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *OGame) ServerTime() (time.Time, error) {
	return b.WithPriority(taskRunner.Normal).ServerTime()
}

// Location returns bot Time zone.
func (b *OGame) Location() *time.Location {
	return b.location
}

// GetUserInfos gets the user information
func (b *OGame) GetUserInfos() (ogame.UserInfos, error) {
	return b.WithPriority(taskRunner.Normal).GetUserInfos()
}

// SendMessage sends a message to playerID
func (b *OGame) SendMessage(playerID int64, message string) error {
	return b.WithPriority(taskRunner.Normal).SendMessage(playerID, message)
}

// SendMessageAlliance sends a message to associationID
func (b *OGame) SendMessageAlliance(associationID int64, message string) error {
	return b.WithPriority(taskRunner.Normal).SendMessageAlliance(associationID, message)
}

// GetFleets get the player's own fleets activities
func (b *OGame) GetFleets(opts ...Option) ([]ogame.Fleet, ogame.Slots) {
	return b.WithPriority(taskRunner.Normal).GetFleets(opts...)
}

// GetFleetsFromEventList get the player's own fleets activities
func (b *OGame) GetFleetsFromEventList() []ogame.Fleet {
	return b.WithPriority(taskRunner.Normal).GetFleetsFromEventList()
}

// CancelFleet cancel a fleet
func (b *OGame) CancelFleet(fleetID ogame.FleetID) error {
	return b.WithPriority(taskRunner.Normal).CancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *OGame) GetAttacks(opts ...Option) ([]ogame.AttackEvent, error) {
	return b.WithPriority(taskRunner.Normal).GetAttacks(opts...)
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *OGame) GalaxyInfos(galaxy, system int64, options ...Option) (ogame.SystemInfos, error) {
	return b.WithPriority(taskRunner.Normal).GalaxyInfos(galaxy, system, options...)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *OGame) GetResourceSettings(planetID ogame.PlanetID, options ...Option) (ogame.ResourceSettings, error) {
	return b.WithPriority(taskRunner.Normal).GetResourceSettings(planetID, options...)
}

// SetResourceSettings set the resources settings on a planet
func (b *OGame) SetResourceSettings(planetID ogame.PlanetID, settings ogame.ResourceSettings) error {
	return b.WithPriority(taskRunner.Normal).SetResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *OGame) GetResourcesBuildings(celestialID ogame.CelestialID, options ...Option) (ogame.ResourcesBuildings, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesBuildings(celestialID, options...)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *OGame) GetDefense(celestialID ogame.CelestialID, options ...Option) (ogame.DefensesInfos, error) {
	return b.WithPriority(taskRunner.Normal).GetDefense(celestialID, options...)
}

// GetShips gets all ships units information of a planet
func (b *OGame) GetShips(celestialID ogame.CelestialID, options ...Option) (ogame.ShipsInfos, error) {
	return b.WithPriority(taskRunner.Normal).GetShips(celestialID, options...)
}

// GetFacilities gets all facilities information of a planet
func (b *OGame) GetFacilities(celestialID ogame.CelestialID, options ...Option) (ogame.Facilities, error) {
	return b.WithPriority(taskRunner.Normal).GetFacilities(celestialID, options...)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *OGame) GetProduction(celestialID ogame.CelestialID) ([]ogame.Quantifiable, int64, error) {
	return b.WithPriority(taskRunner.Normal).GetProduction(celestialID)
}

// GetCachedResearch returns cached researches
func (b *OGame) GetCachedResearch() ogame.Researches {
	return b.WithPriority(taskRunner.Normal).GetCachedResearch()
}

// GetResearch gets the player researches information
func (b *OGame) GetResearch() (ogame.Researches, error) {
	return b.WithPriority(taskRunner.Normal).GetResearch()
}

// GetSlots gets the player current and total slots information
func (b *OGame) GetSlots() (ogame.Slots, error) {
	return b.WithPriority(taskRunner.Normal).GetSlots()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *OGame) Build(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).Build(celestialID, id, nbr)
}

// TechnologyDetails extract details from ajax window when clicking supplies/facilities/techs/lf...
func (b *OGame) TechnologyDetails(celestialID ogame.CelestialID, id ogame.ID) (ogame.TechnologyDetails, error) {
	return b.WithPriority(taskRunner.Normal).TechnologyDetails(celestialID, id)
}

// TearDown tears down any ogame building
func (b *OGame) TearDown(celestialID ogame.CelestialID, id ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).TearDown(celestialID, id)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *OGame) BuildCancelable(celestialID ogame.CelestialID, id ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *OGame) BuildProduction(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *OGame) BuildBuilding(celestialID ogame.CelestialID, buildingID ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(celestialID ogame.CelestialID, defenseID ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(celestialID ogame.CelestialID, shipID ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *OGame) ConstructionsBeingBuilt(celestialID ogame.CelestialID) (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64, ogame.ID, int64) {
	return b.WithPriority(taskRunner.Normal).ConstructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *OGame) CancelBuilding(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelBuilding(celestialID)
}

// CancelLfBuilding cancel the construction of a lifeform building on a specified planet
func (b *OGame) CancelLfBuilding(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelLfBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *OGame) CancelResearch(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelResearch(celestialID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *OGame) BuildTechnology(celestialID ogame.CelestialID, technologyID ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildTechnology(celestialID, technologyID)
}

// GetResources gets user resources
func (b *OGame) GetResources(celestialID ogame.CelestialID) (ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetResources(celestialID)
}

// GetResourcesDetails gets user resources
func (b *OGame) GetResourcesDetails(celestialID ogame.CelestialID) (ogame.ResourcesDetails, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesDetails(celestialID)
}

// GetTechs gets a celestial supplies/facilities/ships/researches
func (b *OGame) GetTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	return b.WithPriority(taskRunner.Normal).GetTechs(celestialID)
}

// SendFleet sends a fleet
func (b *OGame) SendFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).SendFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (b *OGame) EnsureFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).EnsureFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID)
}

// DestroyRockets destroys anti-ballistic & inter-planetary missiles
func (b *OGame) DestroyRockets(planetID ogame.PlanetID, abm, ipm int64) error {
	return b.WithPriority(taskRunner.Normal).DestroyRockets(planetID, abm, ipm)
}

// SendIPM sends IPM
func (b *OGame) SendIPM(planetID ogame.PlanetID, coord ogame.Coordinate, nbr int64, priority ogame.ID) (int64, error) {
	return b.WithPriority(taskRunner.Normal).SendIPM(planetID, coord, nbr, priority)
}

// GetCombatReportSummaryFor gets the latest combat report for a given coordinate
func (b *OGame) GetCombatReportSummaryFor(coord ogame.Coordinate) (ogame.CombatReportSummary, error) {
	return b.WithPriority(taskRunner.Normal).GetCombatReportSummaryFor(coord)
}

// GetEspionageReportFor gets the latest espionage report for a given coordinate
func (b *OGame) GetEspionageReportFor(coord ogame.Coordinate) (ogame.EspionageReport, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReportFor(coord)
}

// GetExpeditionMessages gets the expedition messages
func (b *OGame) GetExpeditionMessages(maxPage int64) ([]ogame.ExpeditionMessage, error) {
	return b.WithPriority(taskRunner.Normal).GetExpeditionMessages(maxPage)
}

// GetExpeditionMessageAt gets the expedition message for time t
func (b *OGame) GetExpeditionMessageAt(t time.Time) (ogame.ExpeditionMessage, error) {
	return b.WithPriority(taskRunner.Normal).GetExpeditionMessageAt(t)
}

// CollectAllMarketplaceMessages collect all marketplace messages
func (b *OGame) CollectAllMarketplaceMessages() error {
	return b.WithPriority(taskRunner.Normal).CollectAllMarketplaceMessages()
}

// CollectMarketplaceMessage collect marketplace message
func (b *OGame) CollectMarketplaceMessage(msg ogame.MarketplaceMessage) error {
	return b.WithPriority(taskRunner.Normal).CollectMarketplaceMessage(msg)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *OGame) GetEspionageReportMessages(maxPage int64) ([]ogame.EspionageReportSummary, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReportMessages(maxPage)
}

// GetEspionageReport gets a detailed espionage report
func (b *OGame) GetEspionageReport(msgID int64) (ogame.EspionageReport, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *OGame) DeleteMessage(msgID int64) error {
	return b.WithPriority(taskRunner.Normal).DeleteMessage(msgID)
}

// DeleteAllMessagesFromTab deletes all messages from a tab in the mail box
func (b *OGame) DeleteAllMessagesFromTab(tabID ogame.MessagesTabID) error {
	return b.WithPriority(taskRunner.Normal).DeleteAllMessagesFromTab(tabID)
}

// GetResourcesProductions gets the planet resources production
func (b *OGame) GetResourcesProductions(planetID ogame.PlanetID) (ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesProductions(planetID)
}

// GetResourcesProductionsLight gets the planet resources production
func (b *OGame) GetResourcesProductionsLight(resBuildings ogame.ResourcesBuildings, researches ogame.Researches,
	resSettings ogame.ResourceSettings, temp ogame.Temperature) ogame.Resources {
	return b.WithPriority(taskRunner.Normal).GetResourcesProductionsLight(resBuildings, researches, resSettings, temp)
}

// FlightTime calculate flight time and fuel needed
func (b *OGame) FlightTime(origin, destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, missionID ogame.MissionID) (secs, fuel int64) {
	return b.WithPriority(taskRunner.Normal).FlightTime(origin, destination, speed, ships, missionID)
}

// Distance return distance between two coordinates
func (b *OGame) Distance(origin, destination ogame.Coordinate) int64 {
	return Distance(origin, destination, b.serverData.Galaxies, b.serverData.Systems, b.serverData.DonutGalaxy, b.serverData.DonutSystem)
}

// RegisterWSCallback ...
func (b *OGame) RegisterWSCallback(id string, fn func(msg []byte)) {
	b.Lock()
	defer b.Unlock()
	b.wsCallbacks[id] = fn
}

// RemoveWSCallback ...
func (b *OGame) RemoveWSCallback(id string) {
	b.Lock()
	defer b.Unlock()
	delete(b.wsCallbacks, id)
}

// RegisterChatCallback register a callback that is called when chat messages are received
func (b *OGame) RegisterChatCallback(fn func(msg ogame.ChatMsg)) {
	b.chatCallbacks = append(b.chatCallbacks, fn)
}

// RegisterAuctioneerCallback register a callback that is called when auctioneer packets are received
func (b *OGame) RegisterAuctioneerCallback(fn func(packet any)) {
	b.auctioneerCallbacks = append(b.auctioneerCallbacks, fn)
}

// RegisterHTMLInterceptor ...
func (b *OGame) RegisterHTMLInterceptor(fn func(method, url string, params, payload url.Values, pageHTML []byte)) {
	b.interceptorCallbacks = append(b.interceptorCallbacks, fn)
}

// Phalanx scan a coordinate from a moon to get fleets information
// IMPORTANT: My account was instantly banned when I scanned an invalid coordinate.
// IMPORTANT: This function DOES validate that the coordinate is a valid planet in range of phalanx
//
//	and that you have enough deuterium.
func (b *OGame) Phalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).Phalanx(moonID, coord)
}

// UnsafePhalanx same as Phalanx but does not perform any input validation.
func (b *OGame) UnsafePhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).UnsafePhalanx(moonID, coord)
}

// JumpGateDestinations returns available destinations for jump gate.
func (b *OGame) JumpGateDestinations(origin ogame.MoonID) (moonIDs []ogame.MoonID, rechargeCountdown int64, err error) {
	return b.WithPriority(taskRunner.Normal).JumpGateDestinations(origin)
}

// JumpGate sends ships through a jump gate.
func (b *OGame) JumpGate(origin, dest ogame.MoonID, ships ogame.ShipsInfos) (success bool, rechargeCountdown int64, err error) {
	return b.WithPriority(taskRunner.Normal).JumpGate(origin, dest, ships)
}

// BuyOfferOfTheDay buys the offer of the day.
func (b *OGame) BuyOfferOfTheDay() error {
	return b.WithPriority(taskRunner.Normal).BuyOfferOfTheDay()
}

// CreateUnion creates a union
func (b *OGame) CreateUnion(fleet ogame.Fleet, users []string) (int64, error) {
	return b.WithPriority(taskRunner.Normal).CreateUnion(fleet, users)
}

// HeadersForPage gets the headers for a specific ogame page
func (b *OGame) HeadersForPage(url string) (http.Header, error) {
	return b.WithPriority(taskRunner.Normal).HeadersForPage(url)
}

// GetEmpire gets all planets/moons information resources/supplies/facilities/ships/researches
func (b *OGame) GetEmpire(celestialType ogame.CelestialType) ([]ogame.EmpireCelestial, error) {
	return b.WithPriority(taskRunner.Normal).GetEmpire(celestialType)
}

// GetEmpireJSON retrieves JSON from Empire page (Commander only).
func (b *OGame) GetEmpireJSON(celestialType ogame.CelestialType) (any, error) {
	return b.WithPriority(taskRunner.Normal).GetEmpireJSON(celestialType)
}

// CharacterClass returns the bot character class
func (b *OGame) CharacterClass() ogame.CharacterClass {
	return b.characterClass
}

// CountColonies returns colonies count/possible
func (b *OGame) CountColonies() (int64, int64) {
	return b.coloniesCount, b.coloniesPossible
}

// GetAuction ...
func (b *OGame) GetAuction() (ogame.Auction, error) {
	return b.WithPriority(taskRunner.Normal).GetAuction()
}

// DoAuction ...
func (b *OGame) DoAuction(bid map[ogame.CelestialID]ogame.Resources) error {
	return b.WithPriority(taskRunner.Normal).DoAuction(bid)
}

// Highscore ...
func (b *OGame) Highscore(category, typ, page int64) (ogame.Highscore, error) {
	return b.WithPriority(taskRunner.Normal).Highscore(category, typ, page)
}

// GetAllResources gets the resources of all planets and moons
func (b *OGame) GetAllResources() (map[ogame.CelestialID]ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetAllResources()
}

// GetTasks return how many tasks are queued in the heap.
func (b *OGame) GetTasks() taskRunner.TasksOverview {
	return b.getTasks()
}

// GetDMCosts returns fast build with DM information
func (b *OGame) GetDMCosts(celestialID ogame.CelestialID) (ogame.DMCosts, error) {
	return b.WithPriority(taskRunner.Normal).GetDMCosts(celestialID)
}

// UseDM use dark matter to fast build
func (b *OGame) UseDM(typ string, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).UseDM(typ, celestialID)
}

// GetItems get all items information
func (b *OGame) GetItems(celestialID ogame.CelestialID) ([]ogame.Item, error) {
	return b.WithPriority(taskRunner.Normal).GetItems(celestialID)
}

// GetActiveItems ...
func (b *OGame) GetActiveItems(celestialID ogame.CelestialID) ([]ogame.ActiveItem, error) {
	return b.WithPriority(taskRunner.Normal).GetActiveItems(celestialID)
}

// ActivateItem activate an item
func (b *OGame) ActivateItem(ref string, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).ActivateItem(ref, celestialID)
}

// BuyMarketplace buy an item on the marketplace
func (b *OGame) BuyMarketplace(itemID int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).BuyMarketplace(itemID, celestialID)
}

// OfferSellMarketplace sell offer on marketplace
func (b *OGame) OfferSellMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).OfferSellMarketplace(itemID, quantity, priceType, price, priceRange, celestialID)
}

// OfferBuyMarketplace buy offer on marketplace
func (b *OGame) OfferBuyMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).OfferBuyMarketplace(itemID, quantity, priceType, price, priceRange, celestialID)
}

// GetLfBuildings ...
func (b *OGame) GetLfBuildings(celestialID ogame.CelestialID, opts ...Option) (ogame.LfBuildings, error) {
	return b.WithPriority(taskRunner.Normal).GetLfBuildings(celestialID, opts...)
}

// GetLfResearch ...
func (b *OGame) GetLfResearch(celestialID ogame.CelestialID, opts ...Option) (ogame.LfResearches, error) {
	return b.WithPriority(taskRunner.Normal).GetLfResearch(celestialID, opts...)
}

// SendDiscoveryFleet ...
func (b *OGame) SendDiscoveryFleet(celestialID ogame.CelestialID, coord ogame.Coordinate, options ...Option) error {
	return b.WithPriority(taskRunner.Normal).SendDiscoveryFleet(celestialID, coord, options...)
}

// SendDiscoveryFleet2 ...
func (b *OGame) SendDiscoveryFleet2(celestialID ogame.CelestialID, coord ogame.Coordinate, options ...Option) (ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).SendDiscoveryFleet2(celestialID, coord, options...)
}

// GetAvailableDiscoveries ...
func (b *OGame) GetAvailableDiscoveries(opts ...Option) int64 {
	return b.WithPriority(taskRunner.Normal).GetAvailableDiscoveries(opts...)
}

// GetPositionsAvailableForDiscoveryFleet ...
func (b *OGame) GetPositionsAvailableForDiscoveryFleet(galaxy int64, system int64, opts ...Option) ([]ogame.Coordinate, error) {
	return b.WithPriority(taskRunner.Normal).GetPositionsAvailableForDiscoveryFleet(galaxy, system, opts...)
}

// SetProxy this will change the bot http transport object.
// proxyType can be "http" or "socks5".
// An empty proxyAddress will reset the client transport to default value.
func (b *OGame) SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error {
	return b.setProxy(proxyAddress, username, password, proxyType, loginOnly, config)
}
