package gdpr

import (
	tcf2 "github.com/prebid/go-gdpr/vendorconsent/tcf2"
)

type NoEnforcementBuilder func()

// NoEnforcement determines legal basis for a given purpose using defaults
// NoEnforcement implements the PurposeEnforcer interface
type NoEnforcement struct {
	purposeCfg purposeConfig
}

func NewNoEnforcement(cfg purposeConfig) *NoEnforcement {
	return &NoEnforcement{
		purposeCfg: cfg,
	}
}

func (ne *NoEnforcement) LegalBasis(vendorInfo VendorInfo, bidderInfo BidderInfo, consent tcf2.ConsentMetadata) bool {
	return true
}

func (ne *NoEnforcement) PurposeEnforced() bool {
	if ne.purposeCfg.EnforcePurpose == TCF2FullEnforcement {
		return true
	}
	return false
}
