package ogame

import (
	"errors"
	"fmt"
	"time"
)

// ErrNotLogged returned when the bot is not logged
var ErrNotLogged = errors.New("not logged")

// ErrMobileView returned when the bot is in mobile view
var ErrMobileView = errors.New("mobile view not supported")

// ErrInvalidPlanetID returned when a planet id is invalid
var ErrInvalidPlanetID = errors.New("invalid planet id")

// ErrAllSlotsInUse returned when all slots are in use
var ErrAllSlotsInUse = errors.New("all slots are in use")

// ErrBotInactive returned when the bot is not active
var ErrBotInactive = errors.New("bot is not active")

// ErrBotLoggedOut returned when the bot is logged out (manually logged out)
var ErrBotLoggedOut = errors.New("bot is logged out")

// ErrDeactivateHidePictures returned when "Hide pictures in reports" is activated
var ErrDeactivateHidePictures = errors.New("deactivate 'Hide pictures in reports'")

// ErrEventsBoxNotDisplayed returned when trying to get attacks from a full page without event box
var ErrEventsBoxNotDisplayed = errors.New("eventList box is not displayed")

// ErrNotEnoughDeuterium ...
var ErrNotEnoughDeuterium = errors.New("not enough deuterium")

// Send fleet errors
var (
	ErrUnionNotFound                      = errors.New("union not found")
	ErrAccountInVacationMode              = errors.New("account in vacation mode")
	ErrNoShipSelected                     = errors.New("no ships to send")
	ErrUninhabitedPlanet                  = errors.New("uninhabited planet")
	ErrNoDebrisField                      = errors.New("no debris field")
	ErrPlayerInVacationMode               = errors.New("player in vacation mode")
	ErrAdminOrGM                          = errors.New("admin or GM")
	ErrNoAstrophysics                     = errors.New("you have to research Astrophysics first")
	ErrNoobProtection                     = errors.New("noob protection")
	ErrPlayerTooStrong                    = errors.New("this planet can not be attacked as the player is to strong")
	ErrNoMoonAvailable                    = errors.New("no moon available")
	ErrNoRecyclerAvailable                = errors.New("no recycler available")
	ErrNoEventsRunning                    = errors.New("there are currently no events running")
	ErrPlanetAlreadyReservedForRelocation = errors.New("this planet has already been reserved for a relocation")
	ErrNotEnoughFuel                      = errors.New("not enough fuel")                          // 4028 Not enough fuel
	ErrAttackBannedUntil                  = errors.New("attack ban until")                         // 4050 Attack ban until
	ErrNotEnoughCargoSpace                = errors.New("not enough cargo space")                   // 140028 Not enough cargo space
	ErrNotEnoughShips                     = errors.New("not enough ships to send")                 // 140054 No ships available
	ErrEngagedInCombat                    = errors.New("the fleet is currently engaged in combat") // 140068 The fleet is currently engaged in combat
)

// AttackBlockActivatedErr ...
type AttackBlockActivatedErr struct {
	BlockedUntil time.Time
}

// Error ...
func (a *AttackBlockActivatedErr) Error() string {
	return fmt.Sprintf("account block activated : %s", a.BlockedUntil.Format(time.DateTime))
}

// NewAttackBlockActivatedErr ...
func NewAttackBlockActivatedErr(blockedUntil time.Time) *AttackBlockActivatedErr {
	return &AttackBlockActivatedErr{BlockedUntil: blockedUntil}
}
