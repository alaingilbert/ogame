package ogame

// AuctioneerNewBid ...
type AuctioneerNewBid struct {
	Sum       int64
	Price     int64
	Bids      int64
	AuctionID int64
	Player    struct {
		ID   int64
		Name string
		Link string
	}
}

// AuctioneerNewAuction ...
// 5::/auctioneer:{"name":"new auction","args":[{"info":"<span style=\"color:#99CC00;\"><b>approx. 45m</b></span> remaining until the auction ends","item":{"uuid":"118d34e685b5d1472267696d1010a393a59aed03","image":"bdb4508609de1df58bf4a6108fff73078c89f777","rarity":"rare"},"oldAuction":{"item":{"uuid":"8a4f9e8309e1078f7f5ced47d558d30ae15b4a1b","imageSmall":"014827f6d1d5b78b1edd0d4476db05639e7d9367","rarity":"rare"},"time":"06.01.2021 17:35:05","bids":1,"sum":1000,"player":{"id":111106,"name":"Governor Skat","link":"http://s152-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=1&system=218"}},"auctionId":18550}]}
type AuctioneerNewAuction struct {
	AuctionID int64
	Approx    int64
}

// AuctioneerAuctionFinished ...
// 5::/auctioneer:{"name":"auction finished","args":[{"sum":2000,"player":{"id":106734,"name":"Someone","link":"http://s152-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=4&system=116"},"bids":2,"info":"Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">1390</span>","time":"06:36"}]}
type AuctioneerAuctionFinished struct {
	Sum         int64
	Bids        int64
	NextAuction int64
	Time        string
	Player      struct {
		ID   int64
		Name string
		Link string
	}
}

// AuctioneerTimeRemaining ...
// 5::/auctioneer:{"name":"timeLeft","args":["<span style=\"color:#FFA500;\"><b>approx. 10m</b></span> remaining until the auction ends"]} // every minute
type AuctioneerTimeRemaining struct {
	Approx int64
}

// AuctioneerNextAuction ...
// 5::/auctioneer:{"name":"timeLeft","args":["Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">598</span>"]}
type AuctioneerNextAuction struct {
	Secs int64
}
