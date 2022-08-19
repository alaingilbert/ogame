package ogame

import "github.com/alaingilbert/ogame/pkg/utils"

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
