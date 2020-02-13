package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vsFactoryAPIEngine/controllers"
	"vsFactoryAPIEngine/core"
)

var apiController *controllers.MongoDBApiController
func init(){
	var err error
	// apiController, err = controllers.NewMongoDBApiController(service.NewMongoService(), "mongodb://127.0.0.1:27017", "test")
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	r := gin.Default()
	// 服务注册
	serviceGroup := r.Group("service"); {
		bundleManager := core.NewRPCBundleManager()
		serviceGroup.POST("/bind", func(c *gin.Context) {
			err = bundleManager.Bind(c)
			if err != nil {
				panic(err)
			}
			// 注册服务
			// 动态添加服务
			err := bundleManager.RPCServe(serviceGroup, "HelloService", core.M{"data":"world1;world2"},
				func(c *gin.Context, request core.IRPCBundleRequest, reply core.IRPCBundleResponse) {
					c.String(http.StatusOK, reply.String())
				})
			if err != nil {
				panic(err)
			}
		})
	}
	// API 服务
	apiGroup := r.Group("api"); {
		apiGroup.POST("/save", apiController.Saver("users"))
		apiGroup.POST("/find", apiController.Finder("users"))
	}
	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
