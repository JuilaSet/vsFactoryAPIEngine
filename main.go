package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vsFactoryAPIEngine/controllers"
	"vsFactoryAPIEngine/serviceRegister"
	"vsFactoryAPIEngine/util"
)

//const username = "admin"

var mongoDBAPIController *controllers.MongoDBApiController
var loginController *controllers.ApplyForServiceController

func init() {
	// api服务控制器
	mongoDBAPIController = controllers.NewMongoDBApiController(serviceRegister.GetMemServiceRegister())
	// 登录服务: 装饰器
	loginController = controllers.NewLoginController(serviceRegister.GetMemServiceRegister())
	//sf := service.GetServiceFactory()
	//sr := serviceRegister.GetMemServiceRegister()
	//s := sf.GetServerByName(service.MongodbService, username, "mongodb://127.0.0.1:27017")
	//sr.Register(username, s)
}

// @title vsFactoryAPI引擎
// @version 0.0.1
// @description  vsFactoryAPI引擎 - 测试
// @BasePath /api/v1/
func main() {
	var err error
	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.GET("/", util.Cors(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
	}))
	// API 服务
	apiGroup := r.Group("api/v1/")
	{
		dbGroup := apiGroup.Group("db")
		{
			// JSON 路由
			dbGroup.Any("/save", util.Cors(mongoDBAPIController.Saver()))
			dbGroup.Any("/find", util.Cors(mongoDBAPIController.Finder()))
			dbGroup.Any("/remove", util.Cors(mongoDBAPIController.Deleter()))
			dbGroup.Any("/update", util.Cors(mongoDBAPIController.Updater()))
		}
		// Restful 路由
		restfulGroup := apiGroup.Group("mongodb-service"); {
			// JSON 路由
			restfulGroup.POST("/:"+controllers.RestfulCollectionName, util.Cors(mongoDBAPIController.RestfulSaver()))
			restfulGroup.GET("/:"+controllers.RestfulCollectionName, util.Cors(mongoDBAPIController.RestfulFinder()))
			restfulGroup.DELETE("/:"+controllers.RestfulCollectionName, util.Cors(mongoDBAPIController.RestfulDeleter()))
			restfulGroup.PUT("/:"+controllers.RestfulCollectionName, util.Cors(mongoDBAPIController.RestfulUpdater()))
		}
		// 验证 服务
		sessGroup := apiGroup.Group("sess"); {
			// 登录 {servername, serverURL}	// "mongodb://127.0.0.1:27017"
			sessGroup.Any("/login", util.Cors(loginController.ApplyForServiceHandler()))
			// 查寻登录情况 {servername} // service.MongodbService
			sessGroup.Any("/loginCheck", util.Cors(loginController.LoginCheckHandler()))
		}
	}

	// 文档界面访问URL
	// http://127.0.0.1:8080/swagger/index.html
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
