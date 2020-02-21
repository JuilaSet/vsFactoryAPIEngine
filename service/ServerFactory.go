package service

import "errors"

type Server interface{}

type ServerFactory struct{}

// 单例模式
var serverFactory *ServerFactory

func init() {
	serverFactory = &ServerFactory{}
}

func GetServiceFactory() *ServerFactory {
	return serverFactory
}

// 服务对象字符串
const (
	MongodbService = "mongodb-service"
)

// "mongodb-service" -> mongodb数据库服务
func (s *ServerFactory) GetServerByName(servername string, userID string, serverURL string) Server {
	switch servername {
	case MongodbService:
		s, err := NewMongoService(serverURL, userID)
		if err != nil {
			panic(err)
		}
		return s
	default:
		panic(errors.New("server not support"))
	}
}
