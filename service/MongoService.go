package service

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"vsFactoryAPIEngine/dataAdapter"
)

// implement IJsonDBService
type MongoService struct {
	sess *mgo.Session
	db   *mgo.Database
}

func NewMongoService() *MongoService {
	return &MongoService{}
}

func (m *MongoService) Connect(url string, dbName string) (err error) {
	sess, err := mgo.Dial(url)
	db := sess.DB(dbName)
	m.sess = sess
	m.db = db
	return
}

func (m *MongoService) Save(collectionName string, data dataAdapter.IJsonAdapter) (err error) {
	collection := m.db.C(collectionName)
	d, err := data.Data()
	if err != nil {
		return
	}
	err = collection.Insert(d)
	return
}

func (m *MongoService) Find(collectionName string, query dataAdapter.IJsonAdapter) (results []bson.M, err error) {
	collection := m.db.C(collectionName)
	q, err := query.Data()
	err = collection.Find(q).All(&results)
	return
}
