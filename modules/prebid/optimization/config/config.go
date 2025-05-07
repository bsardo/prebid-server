package config

import (
	"encoding/json"
	"fmt"

	"github.com/prebid/prebid-server/v3/util/jsonutil"
)

func NewConfig(data json.RawMessage) (*PbRulesEngine, error) {
	cfg := &PbRulesEngine{}

	if err := jsonutil.UnmarshalValid(data, cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
