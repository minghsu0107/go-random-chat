package match

import (
	"encoding/json"
)

type MatchResultPresenter struct {
	AccessToken string `json:"access_token"`
}

func (m *MatchResultPresenter) Encode() []byte {
	result, _ := json.Marshal(m)
	return result
}
