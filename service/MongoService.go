package service

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// implement IMongoDBService
type MongoService struct {
	sess *mgo.Session
	db   *mgo.Database
}

func NewMongoService(url string, dbName string) (*MongoService, error) {
	sess, err := mgo.Dial(url)
	db := sess.DB(dbName)
	m := &MongoService{}
	m.sess = sess
	m.db = db
	return m, err
}

func (m *MongoService) Save(collectionName string, data bson.M) (err error) {
	collection := m.db.C(collectionName)
	err = collection.Insert(data)
	return
}

func (m *MongoService) Find(collectionName string, query bson.M) (results []bson.M, err error) {
	collection := m.db.C(collectionName)
	if str, ok := query["_id"].(string); ok {
		query["_id"] = bson.ObjectIdHex(str)
	}
	err = collection.Find(query).All(&results)
	return
}

func (m *MongoService) FindID(collectionName string, id string) (result bson.M, err error) {
	collection := m.db.C(collectionName)
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func (m *MongoService) Delete(collectionName string, query bson.M) (err error) {
	collection := m.db.C(collectionName)
	if str, ok := query["_id"].(string); ok {
		query["_id"] = bson.ObjectIdHex(str)
	}
	err = collection.Remove(query)
	return
}

func (m *MongoService) DeleteID(collectionName string, id string) (err error) {
	collection := m.db.C(collectionName)
	err = collection.RemoveId(bson.ObjectIdHex(id))
	return
}

func (m *MongoService) Update(collectionName string, query bson.M, update bson.M) (err error) {
	collection := m.db.C(collectionName)
	if str, ok := query["_id"].(string); ok {
		query["_id"] = bson.ObjectIdHex(str)
	}
	err = collection.Update(query, update)
	return
}

func (m *MongoService) UpdateID(collectionName string, id string, update bson.M) (err error) {
	collection := m.db.C(collectionName)
	err = collection.UpdateId(bson.ObjectIdHex(id), update)
	return
}

// 断开连接
func (m *MongoService) Disconnect() {
	m.sess.Close()
}
