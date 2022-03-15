package match

import (
	"encoding/json"
)

type User struct {
	ID   uint64
	Name string
}

type MatchResult struct {
	Matched     bool
	UserID      uint64
	PeerID      uint64
	ChannelID   uint64
	AccessToken string
}

func (r *MatchResult) Encode() []byte {
	result, _ := json.Marshal(r)
	return result
}
func (r *MatchResult) ToPresenter() *MatchResultPresenter {
	return &MatchResultPresenter{
		AccessToken: r.AccessToken,
	}
}
