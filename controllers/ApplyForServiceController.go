package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"vsFactoryAPIEngine/Models"
	"vsFactoryAPIEngine/service"
	"vsFactoryAPIEngine/serviceRegister"
)

// 用户申请服务控制器
type ApplyForServiceController struct {
	register serviceRegister.IServiceRegister	// 单例对象
	factory *service.ServerFactory	// 单例对象
}

func NewLoginController(register serviceRegister.IServiceRegister) *ApplyForServiceController {
	factory := service.GetServiceFactory()
	return &ApplyForServiceController{register, factory}
}

// 登录 {servername, serverURL}	// "mongodb://127.0.0.1:27017"
func (controller *ApplyForServiceController) ApplyForServiceHandler(domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		userLoginJson := models.UserApplyForServiceInfo{}
		if err := c.ShouldBindJSON(&userLoginJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "ok" : false})
			return
		}

		// 判断合法性
		// uuid
		uuidString, err := c.Cookie("uuid")
		if err != nil {
			// 没有获取到UUID, 发送UUID给客户端Cookie
			// 并不能通过这个判断用户登录成功
			u1 := uuid.NewV4()
			uuidString := u1.String()
			c.SetCookie("uuid", uuidString,3600,"/",domain,false,true)
		}

		// 查询是否注册了该服务, 防止重复注册: (连接池有服务)
		if controller.register.Valid(uuidString, userLoginJson.ServerName) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already login", "ok" : false})
			return
		}

		// 注册服务
		s := controller.factory.GetServerByName(userLoginJson.ServerName, uuidString, userLoginJson.ServerURL)
		controller.register.Register(uuidString, userLoginJson.ServerName, s)

		// 回显
		c.JSON(http.StatusBadRequest, gin.H{"ok": true})
	}
}

// 查寻登录情况 {servername}
func (controller *ApplyForServiceController) LoginCheckHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()
		userLoginJson := models.UserApplyForServiceInfo{}
		if err := c.ShouldBindJSON(&userLoginJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// uuid
		uuidString, err := c.Cookie("uuid")
		if err != nil {
			// 未进行登录
			c.String(http.StatusOK, "false")
			return
		}

		// 查询是否注册了该服务
		if controller.register.Valid(uuidString, userLoginJson.ServerName) {
			c.String(http.StatusOK, "true")
		}else{
			c.String(http.StatusOK, "false")
		}
	}
}
