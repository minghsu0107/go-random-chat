package match

import (
	"encoding/json"
)

func DecodeToMatchResult(data []byte) (*MatchResult, error) {
	var result MatchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
