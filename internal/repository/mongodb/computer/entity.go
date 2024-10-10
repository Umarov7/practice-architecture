package computer

import "go.mongodb.org/mongo-driver/bson/primitive"

type Computer struct {
	ID           *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IP           string              `json:"ip" bson:"ip"`
	Manufacturer string              `json:"manufacturer" bson:"manufacturer"`
	CPU          string              `json:"cpu" bson:"cpu"`
	RAM          string              `json:"ram" bson:"ram"`
	HDD          string              `json:"hdd" bson:"hdd"`
	GPU          string              `json:"gpu" bson:"gpu"`
	OS           string              `json:"os" bson:"os"`
	IsDeleted    bool                `json:"isDeleted" bson:"isDeleted"`
}
