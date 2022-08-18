package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
)

type Options struct {
	DebugGalaxy     bool
	SkipInterceptor bool
	SkipRetry       bool
	ChangePlanet    ogame.CelestialID // cp parameter
}

// Option functions to be passed to public interface to change behaviors
type Option func(*Options)

// DebugGalaxy option to debug galaxy
func DebugGalaxy(opt *Options) {
	opt.DebugGalaxy = true
}

// SkipInterceptor option to skip html interceptors
func SkipInterceptor(opt *Options) {
	opt.SkipInterceptor = true
}

// SkipRetry option to skip retry
func SkipRetry(opt *Options) {
	opt.SkipRetry = true
}

// ChangePlanet set the cp parameter
func ChangePlanet(celestialID ogame.CelestialID) Option {
	return func(opt *Options) {
		opt.ChangePlanet = celestialID
	}
}

// ChatPayload ...
type ChatPayload struct {
	Name string    `json:"name"`
	Args []ChatMsg `json:"args"`
}

// ChatMsg ...
type ChatMsg struct {
	SenderID      int64  `json:"senderId"`
	SenderName    string `json:"senderName"`
	AssociationID int64  `json:"associationId"`
	Text          string `json:"text"`
	ID            int64  `json:"id"`
	Date          int64  `json:"date"`
}

func (m ChatMsg) String() string {
	return "\n" +
		"     Sender ID: " + utils.FI64(m.SenderID) + "\n" +
		"   Sender name: " + m.SenderName + "\n" +
		"Association ID: " + utils.FI64(m.AssociationID) + "\n" +
		"          Text: " + m.Text + "\n" +
		"            ID: " + utils.FI64(m.ID) + "\n" +
		"          Date: " + utils.FI64(m.Date)
}
