package model

import "encoding/json"

type Match struct {
	RecipientUserID string `json:"recipientUserID" bson:"recipientUserID"`
	ActorUserID     string `json:"actorUserID" bson:"actorUserID"`
	Liked           bool   `json:"liked" bson:"liked"`
	Matched         bool   `json:"matched" bson:"matched"`
}

func (m *Match) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
