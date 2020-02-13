package dataAdapter

import "github.com/globalsign/mgo/bson"

type IJsonAdapter interface {
	Data() (bson.M, error)
	SetData(interface{}) error
}
