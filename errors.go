package ogame

import "errors"

// ErrNotLogged returned when the bot is not logged
var ErrNotLogged = errors.New("not logged")

// ErrBadCredentials returned when the provided credentials are invalid
var ErrBadCredentials = errors.New("bad credentials")

// ErrNotLogged returned when the account is not found
var ErrAccountNotFound = errors.New("account not found")

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

// Send fleet errors
var (
	ErrNoShipSelected                     = errors.New("no ships to send")
	ErrNotEnoughShips                     = errors.New("not enough ships to send")
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
	ErrPlanetAlreadyReservecForRelocation = errors.New("this planet has already been reserved for a relocation")
)
