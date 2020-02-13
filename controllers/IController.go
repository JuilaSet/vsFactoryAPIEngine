package controllers

import (
	"github.com/gin-gonic/gin"
)

type IDBApiController interface {
	Saver(collectionName string) func(c *gin.Context)
	Finder(collectionName string) func(c *gin.Context)
	//Remover() func(c *gin.Context)
}
