package service

import (
	"github.com/globalsign/mgo/bson"
)

type IMongoDBService interface {
	Save(collectionName string, data bson.M) (err error)
	Find(collectionName string, query bson.M) (results []bson.M, err error)
	Delete(collectionName string, query bson.M) (err error)
	//Update(collectionName string, queryAndUpdate bson.M) (err error)
	Disconnect()
}
