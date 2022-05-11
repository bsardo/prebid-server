package gdpr

import (
	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/consentconstants"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type BasicEnforcementBuilder func()()
// type LegalBasisCalculator func(VendorInfo, BidderInfo) (bool, error)

// BasicEnforcement determines legal basis for a given purpose using the TCF2 basic enforcement algorithm
// BasicEnforcement implements the PurposeEnforcer interface
type BasicEnforcement struct {
	cfg purposeConfig
	// legalBasisCalculator LegalBasisCalculator
}

// func NewBasicEnforcement(cfg purposeConfig) (*BasicEnforcement) {
func NewBasicEnforcement(cfg purposeConfig) (PurposeEnforcer) {
	if cfg.purposeID == consentconstants.Purpose(2) {
		return &Purpose2BasicEnforcement{cfg: cfg}
	}
	return &BasicEnforcement{}
}

func (be *BasicEnforcement) LegalBasis(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	if _, found := be.cfg.VendorExceptionMap[bidderInfo.bidder]; found {
		return true, nil
	}

	

	// be.cfg.enforcePurpose
	return false, nil //TODO
}

// func (be *BasicEnforcement) LegalBasisForHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {

// }


func (be *BasicEnforcement) AllowBidRequest(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) { // * THE BIDDER IS THE VENDOR SO WHY NOT PACKAGE THAT INFO TOGETHER SO JUST ONE PARAMETER???
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}
func (be *BasicEnforcement) AllowGeoData(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}
func (be *BasicEnforcement) AllowUserID(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}
func (be *BasicEnforcement) AllowAnalytics(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}
func (be *BasicEnforcement) AllowHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}
func (be *BasicEnforcement) AllowBidderSync(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil	//TODO - do we need this function? Use LegalBasis instead?
}


type Purpose2BasicEnforcement struct {
	cfg purposeConfig
}

func (p2be *Purpose2BasicEnforcement) LegalBasis(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) {
	return false, nil //TODO
}