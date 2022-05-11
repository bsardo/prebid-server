package gdpr

type AuctionPermissions struct {
	allowAnalytics  bool
	allowBidRequest bool
	passGeo         bool
	passID          bool
}

var AllowAll = AuctionPermissions{
	allowBidRequest: true,
	passGeo:         true,
	passID:          true,
}

var DenyAll = AuctionPermissions{
	allowBidRequest: false,
	passGeo:         false,
	passID:          false,
}

var AllowBidRequestOnly = AuctionPermissions {
	allowBidRequest: true,
	passGeo:         false,
	passID:          false,
}