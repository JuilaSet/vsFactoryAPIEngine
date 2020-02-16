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
	err = collection.Find(query).All(&results)
	return
}

func (m *MongoService) Delete(collectionName string, query bson.M) (err error) {
	collection := m.db.C(collectionName)
	err = collection.Remove(query)
	return
}

// 更新: json对象必须包含 update 和 query
//func (m *MongoService) Update(collectionName string, queryAndUpdate dataAdapter.IJsonAdapter) (err error) {
//	collection := m.db.C(collectionName)
//	queryJson, err := queryAndUpdate.QueryAndUpdateData()
//	if err != nil {
//		return
//	}
//	err = collection.Update(queryJson.Query, queryJson.Update)
//	return
//}

// 断开连接
func (m *MongoService) Disconnect(){
	m.sess.Close()
}