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
	register serviceRegister.IServiceRegister // 单例对象
	factory  *service.ServerFactory           // 单例对象
}

func NewLoginController(register serviceRegister.IServiceRegister) *ApplyForServiceController {
	factory := service.GetServiceFactory()
	return &ApplyForServiceController{register, factory}
}

// 登录 {servername, serverURL}	// "mongodb://127.0.0.1:27017"
func (controller *ApplyForServiceController) ApplyForServiceHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		userLoginJson := models.UserApplyForServiceInfo{}
		if err := c.ShouldBindJSON(&userLoginJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "ok": false})
			return
		}

		// uuid
		uuidString := userLoginJson.Username

		// 判断合法性
		if uuidString == "" {
			// 客户端没有获取到UUID, 发送UUID给客户端
			u1 := uuid.NewV4()
			uuidString = u1.String()
			c.JSON(http.StatusOK, gin.H{"ok": true, "username": uuidString})
		}

		// 查询是否注册了该服务, 防止重复注册: (连接池有服务)
		if controller.register.Valid(uuidString, userLoginJson.ServerName) {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "already login"})
			return
		}

		// 注册服务
		s := controller.factory.GetServerByName(userLoginJson.ServerName, uuidString, userLoginJson.ServerURL)
		controller.register.Register(uuidString, userLoginJson.ServerName, s)
		c.JSON(http.StatusOK, gin.H{"ok": true, "username": uuidString})
	}
}

// 查寻登录情况 {servername}
func (controller *ApplyForServiceController) LoginCheckHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		defer service.ControllerRecover()

		// 获取用户登录信息
		userLoginJson := models.UserApplyForServiceInfo{}
		if err := c.ShouldBindJSON(&userLoginJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
			return
		}

		// uuid
		uuidString := userLoginJson.Username
		//uuidString, err := c.Cookie("uuid")
		ok := uuidString != ""
		if !ok {
			// 未进行登录
			c.JSON(http.StatusOK, gin.H{"ok": true, "data": false})
			return
		}

		// 查询是否注册了该服务
		if controller.register.Valid(uuidString, userLoginJson.ServerName) {
			c.JSON(http.StatusOK, gin.H{"ok": true, "data": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"ok": true, "data": false})
		}
	}
}
