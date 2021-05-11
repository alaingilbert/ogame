package ogame

// Item Is an ogame item that can be activated
type Item struct {
	Ref            string
	Name           string
	Image          string
	ImageLarge     string
	Title          string
	Rarity         string // common
	Amount         int64
	AmountFree     int64
	AmountBought   int64
	canBeActivated bool
	//Category                []string
	//Currency                string // dm
	//Costs                   int64
	//IsReduced               bool
	//buyable                 bool
	//canBeBoughtAndActivated bool
	//isAnUpgrade             bool
	//isCharacterClassItem    bool
	//hasEnoughCurrency       bool
	//Cooldown                bool
	//extendable              bool
	//MoonOnlyItem            bool
	//duration                interface{}
	//DurationExtension       interface{}
	//TotalTime               interface{}
	//timeLeft                interface{}
	//status                  interface{}
	//firstStatus             string // effecting
	//ToolTip                 string
	//buyTitle                string
	//activationTitle         string
}

// ActiveItem ...
type ActiveItem struct {
	ID            int64
	Ref           string
	Name          string
	TimeRemaining int64
	TotalDuration int64
	ImgSmall      string
}
