package serviceRegister

import "vsFactoryAPIEngine/service"

// 服务注册对象
type IServiceRegister interface {
	// 用户申请服务 (服务名称 -> 用户ID), 并向对象池中放入服务对象
	Register(uuid string, serviceName string, serve service.Server)
	// 获取用户的服务池
	GetPool(uid string) (IUserServicePool, bool)
	// 判断这个服务是否能够被用户使用
	Valid(uid string, servername string) bool
}

// 用户服务结构体: 存放用户的服务以及用户注册信息
type IUserServicePool interface {
	// 注册服务: 服务名称 -> service, name map[string]Server
	RegisterServe(servername string, server service.Server)
	// 获取服务
	GetServe(servername string) (service.Server, bool)
}

