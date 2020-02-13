package dataAdapter

import (
	"encoding/json"
	"errors"
	"github.com/globalsign/mgo/bson"
)

// implement IJsonAdapter
type MongoDBAdapter struct {
	data interface{}
}

// data must be string, []byte or bson.M
func NewJsonDBAdapter(data interface{}) *MongoDBAdapter {
	return &MongoDBAdapter{data}
}

func (a *MongoDBAdapter) Data() (bson.M, error) {
	switch v := a.data.(type) {
	case []byte:
		m := make(bson.M)
		err := json.Unmarshal(v, &m)
		return m, err
	case string:
		m := make(bson.M)
		err := json.Unmarshal([]byte(v), &m)
		return m, err
	case bson.M:
		return v, nil
	default:
		return nil, errors.New("unknown type of data: need string, []byte or bson.M")
	}
}

func (a *MongoDBAdapter) SetData(d interface{}) (err error) {
	switch v := a.data.(type) {
	case []byte:
		m := make(bson.M)
		err = json.Unmarshal(v, &m)
		if err != nil {
			return
		}
		a.data = m
	case string:
		m := make(bson.M)
		err = json.Unmarshal([]byte(v), &m)
		if err != nil {
			return
		}
		a.data = m
	case bson.M:
		a.data = v
	default:
		return errors.New("unknown type of data: need string, []byte or bson.M")
	}
	return
}
