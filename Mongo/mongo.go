// Package Mongo. Mongo's models
package Mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection Name in MongoDB.
const (
	ColAnother = "anothermessages"
)

// Create a Session to Connect MongoDB Usage: 'Mongo.S.Clone()'
var S *mgo.Session

type ExampleModel struct {
	Id   bson.ObjectId `json:"_id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}
