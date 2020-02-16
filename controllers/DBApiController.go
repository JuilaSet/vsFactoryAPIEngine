package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"vsFactoryAPIEngine/service"
	"vsFactoryAPIEngine/serviceRegister"
)

type UserInsertInfo struct {
	CollectionName string `json:"collectionName"`
	Data bson.M `json:"data"`
}

// implement IDBApiController
type MongoDBApiController struct {
	register serviceRegister.IServiceRegister	// 单例对象
}

func NewMongoDBApiController(register serviceRegister.IServiceRegister) (mc *MongoDBApiController) {
	return &MongoDBApiController{register}
}

// 返回一个路由函数, 用于将客户端传来的json数据存放到数据库中
// {collectionName, data}
func (controller *MongoDBApiController) Saver() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		// 获取用户的UUID
		uuidStr, err := c.Cookie("uuid")
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			panic(errors.New("service should be applied first"))
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取用户发来的json
		info := UserInsertInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok" : false, "error" : "json 格式错误"})
			return
		}
		data := info.Data
		collectionName := info.CollectionName

		// 获取需要插入的数据
		err = dbService.Save(collectionName, data)
		if err != nil {
			panic(err)
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok" : true})
	}
}

//// 返回一个路由函数, 用于将客户端传来的json数据到数据库中进行查询
//func (controller *MongoDBApiController) Finder(collectionName string) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		defer service.ControllerRecover()
//		data, err := ioutil.ReadAll(c.Request.Body)
//		if err != nil {
//			panic(err)
//		}
//
//		dbService := controller.service
//
//		jsonAdaptor := dataAdapter.NewJsonDBAdapter(data)
//		res, err := dbService.Find(collectionName, jsonAdaptor)
//		if err != nil {
//			panic(err)
//		}
//
//		str := make([]string, 0, 10)
//		for _, r := range res {
//			s, err := json.Marshal(r)
//			if err != nil {
//				panic(err)
//			}
//			str = append(str, string(s))
//		}
//		// 写回客户端
//		c.String(http.StatusOK, "Found: \n"+strings.Join(str, ", "))
//	}
//}
//
//// 返回一个路由函数, 用于将客户端传来的json数据到数据库中进行查询并删除该数据
//func (controller *MongoDBApiController) Deleter(collectionName string) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		defer service.ControllerRecover()
//		data, err := ioutil.ReadAll(c.Request.Body)
//		if err != nil {
//			panic(err)
//		}
//
//		dbService := controller.service
//
//		jsonAdaptor := dataAdapter.NewJsonDBAdapter(data)
//		err = dbService.Delete(collectionName, jsonAdaptor)
//		if err != nil {
//			panic(err)
//		}
//
//		str := make([]string, 0, 10)
//		for _, r := range data {
//			s, err := json.Marshal(r)
//			if err != nil {
//				panic(err)
//			}
//			str = append(str, string(s))
//		}
//
//		// 写回客户端
//		c.String(http.StatusOK, "Delete: \n" + strings.Join(str, ", "))
//	}
//}
//
//// 返回一个路由函数, 用于将客户端传来的json数据存到数据库中进行查询并更新
//func (controller *MongoDBApiController) Updater(collectionName string) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		defer service.ControllerRecover()
//		data, err := ioutil.ReadAll(c.Request.Body)
//		if err != nil {
//			panic(err)
//		}
//		dbService := controller.service
//
//		jsonAdaptor := dataAdapter.NewJsonDBAdapter(data)
//		err = dbService.Update(collectionName, jsonAdaptor)
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{"error" : err.Error()})
//			panic(err)
//		}
//		// 写回客户端
//		c.String(http.StatusOK, "Update: \n" + string(data))
//	}
//}
