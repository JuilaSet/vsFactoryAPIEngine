package service

import (
	"github.com/globalsign/mgo/bson"
	"vsFactoryAPIEngine/dataAdapter"
)

type IJsonDBService interface {
	Connect(url string, dbName string) (err error)
	Save(collectionName string, data dataAdapter.IJsonAdapter) (err error)
	Find(collectionName string, query dataAdapter.IJsonAdapter) (results []bson.M, err error)
}
