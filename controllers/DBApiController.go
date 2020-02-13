package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"vsFactoryAPIEngine/dataAdapter"
	"vsFactoryAPIEngine/service"
)

// implement IDBApiController
type MongoDBApiController struct {
	service service.IJsonDBService
	url     string // "mongodb://127.0.0.1:27017"
	dbName  string // test
}

func NewMongoDBApiController(service service.IJsonDBService, url string, dbName string) (mc *MongoDBApiController, err error) {
	mc = &MongoDBApiController{service: service, url: url, dbName: dbName}
	err = service.Connect(url, dbName)
	return
}

func (controller *MongoDBApiController) Saver(collectionName string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		dbService := controller.service
		err = dbService.Connect(controller.url, controller.dbName)
		if err != nil {
			panic(err)
		}

		jsonAdaptor := dataAdapter.NewJsonDBAdapter(data)
		err = dbService.Save(collectionName, jsonAdaptor)
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, "Insert: \n"+string(data))
	}
}

func (controller *MongoDBApiController) Finder(collectionName string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		dbService := controller.service

		jsonAdaptor := dataAdapter.NewJsonDBAdapter(data)
		res, err := dbService.Find(collectionName, jsonAdaptor)
		if err != nil {
			panic(err)
		}

		str := make([]string, 0, 10)
		for _, r := range res {
			s, err := json.Marshal(r)
			if err != nil {
				panic(err)
			}
			str = append(str, string(s))
		}
		c.String(http.StatusOK, "Found: \n"+strings.Join(str, ", "))
	}
}
