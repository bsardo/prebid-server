package gdpr

type FullEnforcementBuilder func()()

// FullEnforcement determines legal basis for a given purpose using the TCF2 full enforcement algorithm
// FullEnforcement implements the PurposeEnforcer interface
type FullEnforcement struct {
	purposeCfg purposeConfig
}

func NewFullEnforcement(cfg purposeConfig) (*FullEnforcement) {
	return &FullEnforcement{}
}

func (fe *FullEnforcement) LegalBasis(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	
}

func (fe *FullEnforcement) LegalBasisForHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}





func (fe *FullEnforcement) AllowBidRequest(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) { // * THE BIDDER IS THE VENDOR SO WHY NOT PACKAGE THAT INFO TOGETHER SO JUST ONE PARAMETER???

}
func (fe *FullEnforcement) AllowGeoData(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}
func (fe *FullEnforcement) AllowUserID(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}
func (fe *FullEnforcement) AllowAnalytics(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}
func (fe *FullEnforcement) AllowHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}
func (fe *FullEnforcement) AllowBidderSync(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

}