package service

import (
	"github.com/globalsign/mgo/bson"
)

type IMongoDBService interface {
	Save(collectionName string, data bson.M) (err error)

	Find(collectionName string, query bson.M) (results []bson.M, err error)
	FindID(collectionName string, id string) (results bson.M, err error)

	Delete(collectionName string, query bson.M) (err error)
	DeleteID(collectionName string, id string) (err error)

	Update(collectionName string, query bson.M, Update bson.M) (err error)

	UpdateID(collectionName string, id string, update bson.M) (err error)
	Disconnect()
}
