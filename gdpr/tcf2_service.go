package gdpr

import (
	"github.com/prebid/prebid-server/openrtb_ext"
)

type TCF2Service struct {
	cfg                   TCF2ConfigReader		// needed for all
	fetchVendorList       VendorListFetcher		// needed for all
	gdprDefaultValue      string				// needed for all
	hostVendorID          int					// needed for hostcookies only
	nonStandardPublishers map[string]struct{}	// needed for auction only
	vendorIDs             map[openrtb_ext.BidderName]uint16	// needed for biddersync and auction only
	// request data
	consent               string				// needed for all
	gdprSignal            Signal				// needed for all
	publisherID           string				// needed for auction only
	aliasGVLIDs           map[string]uint16		// needed for auction only
	// dependency builders
	basicEnforcementBuilder BasicEnforcementBuilder
	fullEnforcementBuilder  FullEnforcementBuilder
	noEnforcementBuilder    NoEnforcementBuilder
	// cached enforcers
	purposeEnforcers map[consentconstants.Purpose]Enforcer
}

// AllowAuctionActivities implements the Permissions interface
// It determines whether auction activities are allowed for a given bidder
func (ts *TCF2Service) AllowAuctionActivities(ctx context.Context, bidder bidderInfo) (permissions AuctionPermissions, err error) {

	if _, ok := ts.nonStandardPublishers[ts.publisherID]; ok {
		return AllowAll, nil
	}
	if ts.gdprSignal == SignalNo {
		return AllowAll, nil
	}
	if ts.consent == "" && ts.gdprSignal == SignalYes {
		return DenyAll, nil
	}

	// note the bidder here is guaranteed to be enabled
	vendorID, vendorFound := ts.resolveVendorId(bidder.bidderCoreName, bidder.bidderName, ts.aliasGVLIDs)
	_, weakVendorEnforcement := ts.cfg.BasicEnforcementVendors()[bidder]

	if !vendorFound && !weakVendorEnforcement {
		return DenyAll, nil
	}
	
	// if the vendorID was not resolved to a valid ID, it will be zero at this point
	parsedConsent, vendor, err := p.parseVendor(ctx, vendorID, consent)
	if err != nil {
		return DenyAll, err
	}

	// vendor will be nil if not a valid TCF2 consent string
	if vendor == nil {
		if weakVendorEnforcement && parsedConsent.Version() == 2 {
			vendor = vendorTrue{}	//TODO
		} else {
			return DenyAll, nil
		}
	}

	if !ts.cfg.IsEnabled() {
		return AllowBidRequestOnly, nil
	}

	consentMeta, ok := parsedConsent.(tcf2.ConsentMetadata)
	if !ok {
		err = fmt.Errorf("Unable to access TCF2 parsed consent")
		return DenyAll, err
	}

	vendorInfo := VendorInfo{vendorID: vendorID, vendor: vendor}
	permissions := AuctionPermissions{}
	permissions.allowBidRequest = allowBidRequest(bidderCoreName, consentMeta, vendorInfo)
	permissions.passGeo = allowGeo(bidderCoreName, consentMeta, vendor)
	permissions.passID = allowID(bidderCoreName, consentMeta, vendorInfo)

	return permissions, nil
}

func (ts *TCF2Service) allowID(bidder openrtb_ext.BidderName, consentMeta tcf2.ConsentMetadata, vendorInfo VendorInfo) bool {
	for i := 2; i <= 10; i++ {
		purpose := consentconstants.Purpose(i)
		enforcer := ts.getPurposeEnforcer(purpose, bidder)

		if enforcer.AllowUserID(...) {
			return true
		}
	}

	return false
}

func (ts *TCFService) allowGeo(bidder openrtb_ext.BidderName, consentMeta tcf2.ConsentMeta, vendor api.Vendor) bool {
	if !ts.cfg.FeatureOneEnforced() {
		return true
	}
	if ts.cfg.FeatureOneVendorException(bidder) {
		return true
	}

	_, weakVendorEnforcement := ts.cfg.BasicEnforcementVendors()[bidder]
	return consentMeta.SpecialFeatureOptIn(1) && (vendor.SpecialFeature(1) || weakVendorEnforcement)
}

func (ts *TCFService) allowBidRequest(bidder openrtb_ext.BidderName, consentMeta tcf2.ConsentMeta, vendorInfo VendorInfo) bool {
	purpose := consentconstants.Purpose(2)
	enforcer := ts.getPurposeEnforcer(purpose, bidder)
	if enforcer.AllowBidRequest() {	// this function will return true if purpose 2 is NOT enforced
		return true
	}
	return false
}

func (ts *TCF2Service) getPurposeEnforcer(purpose consentconstants.Purpose, bidder openrtb_ext.BidderName) Enforcer {
	// use cached enforcer if already exists
	if enforcer, ok := ts.purposeEnforcers[purpose]; ok {
		return enforcer
	}

	cfg := purposeConfig {
		PurposeID:                  purpose,
		EnforcePurpose:             ts.cfg.PurposeEnforced(purpose),
		EnforceVendors:             ts.cfg.PurposeEnforcingVendors(purpose),
		VendorExceptionMap:         ts.cfg.PurposeVendorExceptions(purpose),
		BasicEnforcementVendorsMap: ts.cfg.BasicEnforcementVendors(),
	}
	downgraded := ts.isDowngraded(cfg.EnforcePurpose, ...)
	enforcer := NewPurposeEnforcer(cfg, downgraded)

	//cache the enforcer
	ts.purposeEnforcers[purpose] := enforcer

	return enforcer
}

func (ts *TCF2Service) AllowHostCookies(ctx context.Context) (bool, error) {
	if ts.gdprSignal == SignalNo {
		return true, nil
	}

	return ts.allowSync(ctx, uint16(ts.hostVendorID), false)
}

func (ts *TCF2Service) AllowBidderSync(ctx context.Context, bidder openrtb_ext.BidderName) (bool, error) {
	if ts.gdprSignal == SignalNo {
		return true, nil
	}
	
	id, ok := p.vendorIDs[bidder]
	if !ok {
		return false, nil
	}

	return ts.allowSync(ctx, id, vendorException)
}

func (p *permissionsImpl) allowSync(ctx context.Context, vendorID uint16, vendorException bool) (bool, error) {
	if ts.consent == "" {
		return false, nil
	}

	parsedConsent, vendor, err := ts.parseVendor(ctx, vendorID, ts.consent) //TODO ts.consent does not need to be a parameter
	if err != nil {
		return false, err
	}

	if vendor == nil {
		return false, nil
	}

	if !ts.cfg.PurposeEnforced(consentconstants.Purpose(1)) {
		return true, nil
	}
	consentMeta, ok := parsedConsent.(tcf2.ConsentMetadata)
	if !ok {
		err := errors.New("Unable to access TCF2 parsed consent")
		return false, err
	}

	if ts.cfg.PurposeOneTreatmentEnabled() && consentMeta.PurposeOneTreatment() {
		return ts.cfg.PurposeOneTreatmentAccessAllowed(), nil
	}

	// enforceVendors := p.cfg.PurposeEnforcingVendors(tcf2ConsentConstants.InfoStorageAccess)
	// return p.checkPurpose(consentMeta, vendor, vendorID, tcf2ConsentConstants.InfoStorageAccess, enforceVendors, vendorException, false), nil

	enforcer := ts.getPurposeEnforcer(consentconstants.Purpose(1))
	return enforcer.AllowSync(consentMeta, vendorInfo, bidderInfo)
}

func (p *permissionsImpl) parseVendor(ctx context.Context, vendorID uint16, consent string) (parsedConsent api.VendorConsents, vendor api.Vendor, err error) {
	parsedConsent, err = vendorconsent.ParseString(consent)
	if err != nil {
		err = &ErrorMalformedConsent{
			Consent: consent,
			Cause:   err,
		}
		return
	}

	version := parsedConsent.Version()
	if version != 2 {
		return
	}

	vendorList, err := p.fetchVendorList(ctx, parsedConsent.VendorListVersion())
	if err != nil {
		return
	}

	vendor = vendorList.Vendor(vendorID)
	return
}