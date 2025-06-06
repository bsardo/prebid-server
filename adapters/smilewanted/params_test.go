package smilewanted

import (
	"encoding/json"
	"testing"

	"github.com/prebid/prebid-server/v3/openrtb_ext"
)

// This file actually intends to test static/bidder-params/smilewanted.json
//
// These also validate the format of the external API: request.imp[i].ext.prebid.bidder.smilewanted

// TestValidParams makes sure that the smilewanted schema accepts all imp.ext fields which we intend to support.
func TestValidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json-schemas. %v", err)
	}

	for _, validParam := range validParams {
		if err := validator.Validate(openrtb_ext.BidderSmileWanted, json.RawMessage(validParam)); err != nil {
			t.Errorf("Schema rejected SmileWanted params: %s", validParam)
		}
	}
}

// TestInvalidParams makes sure that the SmileWanted schema rejects all the imp.ext fields we don't support.
func TestInvalidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json-schemas. %v", err)
	}

	for _, invalidParam := range invalidParams {
		if err := validator.Validate(openrtb_ext.BidderSmileWanted, json.RawMessage(invalidParam)); err == nil {
			t.Errorf("Schema allowed unexpected params: %s", invalidParam)
		}
	}
}

var validParams = []string{
	`{"zoneId": "zone_code"}`,
}

var invalidParams = []string{
	`{"zoneId": 100}`,
	`{"zoneId": true}`,
	`{"zoneId": 123}`,
	`{"zoneID": "1"}`,
	``,
	`null`,
	`true`,
	`9`,
	`1.2`,
	`[]`,
	`{}`,
}
