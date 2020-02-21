package models

import "github.com/globalsign/mgo/bson"

// 用户验证信息, 所有请求都必须带有这个信息
type UserVerifyInfo struct {
	Username string `json:"username",default:""` // UUID
}

// 注册服务请求头信息
type UserApplyForServiceInfo struct {
	UserVerifyInfo
	ServerName string `json:"servername"`
	ServerURL  string `json:"serverURL"`
}

// 用户查询服务
type UserQueryInfo struct {
	UserVerifyInfo
	CollectionName string `json:"collectionName"`
	Query          bson.M `json:"query"`
	ID             string `json:"_id"`
}

type UserInsertInfo struct {
	UserVerifyInfo
	CollectionName string `json:"collectionName"`
	Data           bson.M `json:"data"`
}

type UserUpdateInfo struct {
	UserVerifyInfo
	CollectionName string `json:"collectionName"`
	Query          bson.M `json:"query"`
	Update         bson.M `json:"update"`
}

// Restful 接口
type RestfulSaveInfo struct {
	UserVerifyInfo
	Data bson.M `json:"data"`
}

type RestfulFindInfo struct {
	UserVerifyInfo
	ID   string `json:"_id"`
}

type RestfulUpdateInfo struct {
	UserVerifyInfo
	ID   string `json:"_id"`
	Data bson.M `json:"data"`
}

