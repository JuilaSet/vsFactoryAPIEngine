package main

import (
	"github.com/gin-gonic/gin"
	"vsFactoryAPIEngine/controllers"
	"vsFactoryAPIEngine/serviceRegister"
)

var apiController *controllers.MongoDBApiController
var loginController *controllers.ApplyForServiceController
func init(){
	// api服务控制器
	apiController = controllers.NewMongoDBApiController(serviceRegister.GetMemServiceRegister())
	// 登录服务: 装饰器
	loginController = controllers.NewLoginController(serviceRegister.GetMemServiceRegister())
}

func main() {
	var err error
	r := gin.Default()
	// API 服务
	apiGroup := r.Group("api"); {
		dbGroup := apiGroup.Group("db"); {
			dbGroup.POST("/save", apiController.Saver())
			//apiGroup.POST("/find", apiController.Finder("users"))
			//apiGroup.POST("/delete", apiController.Deleter("users"))
			//apiGroup.POST("/update", apiController.Updater("users"))
		}
		// 验证 服务
		sessGroup := apiGroup.Group("sess"); {
			// 登录 {servername, serverURL}	// "mongodb://127.0.0.1:27017"
			sessGroup.POST("/login", loginController.ApplyForServiceHandler("127.0.0.1"))
			// 查寻登录情况 {servername} // service.MongodbService
			sessGroup.POST("/loginCheck", loginController.LoginCheckHandler())
		}
	}
	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
