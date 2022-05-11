package gdpr

import (
	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/consentconstants"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type PurposeEnforcer interface {
	LegalBasis(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)

	// AllowBidRequest(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error) // * THE BIDDER IS THE VENDOR SO WHY NOT PACKAGE THAT INFO TOGETHER SO JUST ONE PARAMETER???
	// AllowGeoData(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)
	// AllowUserID(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)
	// AllowAnalytics(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)
	// AllowHostCookies(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)
	// AllowBidderSync(vendorInfo VendorInfo, bidderInfo BidderInfo) (bool, error)
}

type BidderInfo struct {
	bidderCoreName openrtb_ext.BidderName
	bidder         openrtb_ext.BidderName
}
type VendorInfo struct {
	vendorID uint16
	vendor   api.Vendor
}

type TCF2Enforcement string

const (
	TCF2BasicEnforcement TCF2Enforcement = "basic"
	TCF2FullEnforcement  TCF2Enforcement = "full"
	TCF2NoEnforcement    TCF2Enforcement = "no"
)

type purposeConfig struct {
	purposeID                  consentconstants.Purpose
	enforcePurpose             TCF2Enforcement
	enforceVendors             bool
	// vendorException            bool
	// basicEnforcementVendor     bool
	VendorExceptionMap         map[openrtb_ext.BidderName]struct{}
	BasicEnforcementVendorsMap map[openrtb_ext.BidderName]struct{} // currently map[string]struct{} in agg cfg, only needed for BASIC enforcement
}

// called from the TCF2Service and might be injected into the service as a builder to improve testability
func NewPurposeEnforcer(cfg purposeConfig, downgraded bool) (PurposeEnforcer) {
	if cfg.enforcePurpose == TCF2FullEnforcement && !downgraded {
		return NewFullEnforcement(cfg)
	} else if cfg.enforcePurpose == TCF2BasicEnforcement {
		return NewBasicEnforcement(cfg)
	} else if downgraded {
		return NewBasicEnforcement(cfg)
	}

	return NewNoEnforcement(cfg)
}