package common

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsBotEcho(id string, ids []primitive.ObjectID) bool {
	flag := false

	for _, d := range ids {
		if d.Hex() == id {
			return true
		}
	}

	return flag
}
