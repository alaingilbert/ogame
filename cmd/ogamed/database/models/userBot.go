package models

type UserBot struct {
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
	BotID  uint
	Bot    Bot `gorm:"foreignKey:BotID"`
}
