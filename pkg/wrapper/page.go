package wrapper

//
//type Page struct {
//	b       *OGame
//	doc     *goquery.Document
//	content []byte
//}
//
//func (p *Page) GetDoc() *goquery.Document {
//	if p.doc == nil {
//		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(p.content))
//		p.doc = doc
//	}
//	return p.doc
//}
//
//type EventListAjaxPage struct{ Page }
//type MissileAttackLayerAjaxPage struct{ Page }
//type FetchTechsAjaxPage struct{ Page }
//type RocketlayerAjaxPage struct{ Page }
//type PhalanxAjaxPage struct{ Page }
//
//type FullPage struct{ Page }
//type OverviewPage struct{ FullPage }
//type PreferencesPage struct{ FullPage }
//type SuppliesPage struct{ FullPage }
//type ResourcesSettingsPage struct{ FullPage }
//
////type FacilitiesPageContent struct{ FullPageContent }
////type TraderOverviewPageContent struct{ FullPageContent }
////type TraderResourcesPageContent struct{ FullPageContent }
//
//type ResearchPage struct{ FullPage }
//type FacilitiesPage struct{ FullPage }
//type ShipyardPage struct{ FullPage }
//type DefensesPage struct{ FullPage }
//
////type FleetDispatchPageContent struct{ FullPageContent }
//
//type MovementPage struct{ FullPage }
//
////type GalaxyPageContent struct{ FullPageContent }
////type AlliancePageContent struct{ FullPageContent }
////type PremiumPageContent struct{ FullPageContent }
////type ShopPageContent struct{ FullPageContent }
////type MessagesPageContent struct{ FullPageContent }
////type ChatPageContent struct{ FullPageContent }
////type CharacterClassSelectionPageContent struct{ FullPageContent }
////type BuddiesPageContent struct{ FullPageContent }
////type HighScorePageContent struct{ FullPageContent }
//
//type IFullPage interface {
//	ExtractOGameSession() string
//	ExtractIsInVacation() bool
//	ExtractPlanets() []v6.ExtractorPlanet
//	ExtractAjaxChatToken() (string, error)
//	ExtractCharacterClass() (ogame.CharacterClass, error)
//	ExtractCommander() bool
//	ExtractAdmiral() bool
//	ExtractEngineer() bool
//	ExtractGeologist() bool
//	ExtractTechnocrat() bool
//	ExtractServerTime() (time.Time, error)
//}
//
//func (p PhalanxAjaxPage) ExtractPhalanx() ([]ogame.Fleet, error) {
//	return p.b.extractor.ExtractPhalanx(p.content)
//}
//
//func (p RocketlayerAjaxPage) ExtractDestroyRockets() (int64, int64, string, error) {
//	return p.b.extractor.ExtractDestroyRockets(p.content)
//}
//
//func (p FullPage) ExtractOGameSession() string {
//	return p.b.extractor.ExtractOGameSessionFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractIsInVacation() bool {
//	return p.b.extractor.ExtractIsInVacationFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractAjaxChatToken() (string, error) {
//	return p.b.extractor.ExtractAjaxChatToken(p.content)
//}
//
//func (p FullPage) ExtractCharacterClass() (ogame.CharacterClass, error) {
//	return p.b.extractor.ExtractCharacterClassFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractCommander() bool {
//	return p.b.extractor.ExtractCommanderFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractAdmiral() bool {
//	return p.b.extractor.ExtractAdmiralFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractEngineer() bool {
//	return p.b.extractor.ExtractEngineerFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractGeologist() bool {
//	return p.b.extractor.ExtractGeologistFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractTechnocrat() bool {
//	return p.b.extractor.ExtractTechnocratFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractServerTime() (time.Time, error) {
//	return p.b.extractor.ExtractServerTimeFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractPlanets() []v6.ExtractorPlanet {
//	return p.b.extractor.ExtractPlanetsFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractPlanet(v any) (v6.ExtractorPlanet, error) {
//	return p.b.extractor.ExtractPlanetFromDoc(p.GetDoc(), v)
//}
//
//func (p FullPage) ExtractMoons() []v6.ExtractorMoon {
//	return p.b.extractor.ExtractMoonsFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractMoon(v any) (v6.ExtractorMoon, error) {
//	return p.b.extractor.ExtractMoonFromDoc(p.GetDoc(), v)
//}
//
//func (p FullPage) ExtractCelestials() ([]ogame.ICelestial, error) {
//	return p.b.extractor.ExtractCelestialsFromDoc(p.GetDoc())
//}
//
//func (p FullPage) ExtractCelestial(v any) (ogame.ICelestial, error) {
//	return p.b.extractor.ExtractCelestialFromDoc(p.GetDoc(), v)
//}
//
//func (p ResearchPage) ExtractResearch() ogame.Researches {
//	return p.b.extractor.ExtractResearchFromDoc(p.GetDoc())
//}
//
//func (p SuppliesPage) ExtractResourcesBuildings() (ogame.ResourcesBuildings, error) {
//	return p.b.extractor.ExtractResourcesBuildingsFromDoc(p.GetDoc())
//}
//
//func (p DefensesPage) ExtractDefense() (ogame.DefensesInfos, error) {
//	return p.b.extractor.ExtractDefenseFromDoc(p.GetDoc())
//}
//
//func (p OverviewPage) ExtractActiveItems() ([]ogame.ActiveItem, error) {
//	return p.b.extractor.ExtractActiveItems(p.content)
//}
//
//func (p OverviewPage) ExtractDMCosts() (v6.DMCosts, error) {
//	return p.b.extractor.ExtractDMCosts(p.content)
//}
//
//func (p OverviewPage) ExtractConstructions() (ogame.ID, int64, ogame.ID, int64) {
//	return p.b.extractor.ExtractConstructions(p.content)
//}
//
//func (p OverviewPage) ExtractUserInfos() (ogame.UserInfos, error) {
//	return p.b.extractor.ExtractUserInfos(p.content, p.b.language)
//}
//
//func (p OverviewPage) ExtractCancelResearchInfos() (token string, techID, listID int64, err error) {
//	return p.b.extractor.ExtractCancelResearchInfos(p.content)
//}
//
//func (p OverviewPage) ExtractCancelBuildingInfos() (token string, techID, listID int64, err error) {
//	return p.b.extractor.ExtractCancelBuildingInfos(p.content)
//}
//
//func (p FacilitiesPage) ExtractFacilities() (ogame.Facilities, error) {
//	return p.b.extractor.ExtractFacilitiesFromDoc(p.GetDoc())
//}
//
//func (p ShipyardPage) ExtractProduction() ([]ogame.Quantifiable, int64, error) {
//	return p.b.extractor.ExtractProduction(p.content)
//}
//
//func (p ShipyardPage) ExtractShips() (ogame.ShipsInfos, error) {
//	return p.b.extractor.ExtractShipsFromDoc(p.GetDoc())
//}
//
//func (p ResourcesSettingsPage) ExtractResourceSettings() (ogame.ResourceSettings, error) {
//	return p.b.extractor.ExtractResourceSettingsFromDoc(p.GetDoc())
//}
//
//func (p MovementPage) ExtractFleets() []ogame.Fleet {
//	return p.b.extractor.ExtractFleetsFromDoc(p.GetDoc(), p.b.location)
//}
//
//func (p MovementPage) ExtractSlots() ogame.Slots {
//	return p.b.extractor.ExtractSlotsFromDoc(p.GetDoc())
//}
//
//func (p MovementPage) ExtractCancelFleetToken(fleetID ogame.FleetID) (string, error) {
//	return p.b.extractor.ExtractCancelFleetToken(p.content, fleetID)
//}
//
//func (p EventListAjaxPage) ExtractAttacks(ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
//	return p.b.extractor.ExtractAttacksFromDoc(p.GetDoc(), ownCoords)
//}
//
//func (p MissileAttackLayerAjaxPage) ExtractIPM() (int64, int64, string) {
//	return p.b.extractor.ExtractIPMFromDoc(p.GetDoc())
//}
//
//func (p FetchTechsAjaxPage) ExtractTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, error) {
//	return p.b.extractor.ExtractTechs(p.content)
//}
//
//func (p PreferencesPage) ExtractPreferences() ogame.Preferences {
//	return p.b.extractor.ExtractPreferencesFromDoc(p.GetDoc())
//}
//
//type FullPagePages interface {
//	OverviewPage |
//		PreferencesPage |
//		SuppliesPage |
//		ResourcesSettingsPage |
//		FacilitiesPage |
//		//TraderOverviewPageContent |
//		//TraderResourcesPageContent |
//		ResearchPage |
//		ShipyardPage |
//		DefensesPage |
//		//FleetDispatchPageContent |
//		MovementPage
//	//GalaxyPageContent |
//	//AlliancePageContent |
//	//PremiumPageContent |
//	//ShopPageContent |
//	//MessagesPageContent |
//	//ChatPageContent |
//	//CharacterClassSelectionPageContent |
//	//BuddiesPageContent |
//	//HighScorePageContent
//}
//
//type AjaxPagePages interface {
//	EventListAjaxPage |
//		MissileAttackLayerAjaxPage |
//		FetchTechsAjaxPage |
//		RocketlayerAjaxPage |
//		PhalanxAjaxPage
//}
