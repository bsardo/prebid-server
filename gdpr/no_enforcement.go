package gdpr

type NoEnforcementBuilder func()()

// NoEnforcement determines legal basis for a given purpose using defaults
// NoEnforcement implements the PurposeEnforcer interface
type NoEnforcement struct {
	purposeCfg purposeConfig
}

func NewNoEnforcement(cfg purposeConfig) (*NoEnforcement) {
	return &NoEnforcement{}
}

func (ne *NoEnforcement) AllowBidRequest(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) { // * THE BIDDER IS THE VENDOR SO WHY NOT PACKAGE THAT INFO TOGETHER SO JUST ONE PARAMETER???
	return true, nil
}
func (ne *NoEnforcement) AllowGeoData(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return true, nil
}
func (ne *NoEnforcement) AllowUserID(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return true, nil
}
func (ne *NoEnforcement) AllowAnalytics(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return true, nil
}
func (ne *NoEnforcement) AllowHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return true, nil
}
func (ne *NoEnforcement) AllowBidderSync(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return true, nil
}