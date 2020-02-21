package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"vsFactoryAPIEngine/Models"
	"vsFactoryAPIEngine/service"
	"vsFactoryAPIEngine/serviceRegister"
)

// implement IDBApiController
type MongoDBApiController struct {
	register serviceRegister.IServiceRegister // 单例对象
}

func NewMongoDBApiController(register serviceRegister.IServiceRegister) (mc *MongoDBApiController) {
	return &MongoDBApiController{register}
}

// RestFul路由参数
const (
	RestfulCollectionName = "collectionName"
)

// 返回一个路由函数, 用于将客户端传来的json数据存放到数据库中
// {collectionName, data}
func (controller *MongoDBApiController) Saver() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		// 获取用户发来的json
		info := models.UserInsertInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误: {collectionName, data}"})
			return
		}

		// 获取用户的UUID
		// uuidStr, err := c.Cookie("uuid")
		uuidStr := info.Username
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := info.CollectionName
		if err := dbService.Save(collectionName, info.Data); err != nil {
			panic(err)
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
func (controller *MongoDBApiController) RestfulSaver() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		// 获取用户发来的json
		info := models.RestfulSaveInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误: {username, data}"})
			//c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
			return
		}

		// 获取用户的UUID
		uuidStr := info.Username
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := c.Param(RestfulCollectionName)
		if err := dbService.Save(collectionName, info.Data); err != nil {
			panic(err)
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

// 返回一个路由函数, 用于将客户端传来的json数据到数据库中进行查询
// {collectionName, data}
func (controller *MongoDBApiController) Finder() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.UserQueryInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query}"})
			return
		}

		// 获取用户的UUID
		uuidStr := info.Username
		// uuidStr, err := c.Cookie("uuid")
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := info.CollectionName
		result, err := dbService.Find(collectionName, info.Query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "not found"})
			return
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true, "data": result})
	}
}
func (controller *MongoDBApiController) RestfulFinder() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.RestfulFindInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query}"})
			return
		}

		// 获取用户的UUID
		uuidStr := info.Username
		// uuidStr, err := c.Cookie("uuid")
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
			return
		}

		// 获取需要插入的数据
		collectionName := c.Param(RestfulCollectionName)
		result, err := dbService.FindID(collectionName, info.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "item not found"})
			return
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true, "data": result})
	}
}

// 返回一个路由函数, 用于将客户端传来的json数据到数据库中进行查询并删除该数据
// {collectionName, data}
func (controller *MongoDBApiController) Deleter() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.UserQueryInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query}"})
			return
		}

		// 获取用户的UUID
		uuidStr := info.Username
		//uuidStr, err := c.Cookie("uuid")
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		data := info.Query
		collectionName := info.CollectionName
		if err := dbService.Delete(collectionName, data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "item not found"})
			return
		}

		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
func (controller *MongoDBApiController) RestfulDeleter() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.RestfulFindInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query}"})
			return
		}

		// 获取用户的UUID
		uuidStr := info.Username
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := c.Param(RestfulCollectionName)
		id := info.ID
		if err := dbService.DeleteID(collectionName, id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "item not found"})
			return
		}

		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

//// 返回一个路由函数, 用于将客户端传来的json数据存到数据库中进行查询并更新
func (controller *MongoDBApiController) Updater() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.UserUpdateInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query, update}"})
			return
		}

		// 获取用户的UUID
		//uuidStr, err := c.Cookie("uuid")
		uuidStr := info.Username
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := info.CollectionName
		if err := dbService.Update(collectionName, info.Query, info.Update); err != nil {
			panic(err)
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
func (controller *MongoDBApiController) RestfulUpdater() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户发来的json
		info := models.RestfulUpdateInfo{}
		if err := c.ShouldBindJSON(&info); err != nil {
			// json 格式错误
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "json 格式错误 {collectionName, query, update}"})
			return
		}

		// 获取用户的UUID
		//uuidStr, err := c.Cookie("uuid")
		uuidStr := info.Username
		pool, ok := controller.register.GetPool(uuidStr)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "service should be applied first"})
			return
		}
		s, ok := pool.GetServe(service.MongodbService)
		dbService := s.(service.IMongoDBService)
		if !ok {
			panic(errors.New("service.MongodbService not found"))
		}

		// 获取需要插入的数据
		collectionName := c.Param(RestfulCollectionName) // info.RestfulCollectionName
		if err := dbService.UpdateID(collectionName, info.ID, info.Data); err != nil {
			panic(err)
		}
		// 写回客户端
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
