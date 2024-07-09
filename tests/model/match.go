package model

type Match struct {
	RecipientUserID string `json:"recipientUserID" bson:"recipientUserID"`
	ActorUserID     string `json:"actorUserID" bson:"actorUserID"`
	Liked           bool   `json:"liked" bson:"liked"`
	Matched         bool   `json:"matched" bson:"matched"`
}
