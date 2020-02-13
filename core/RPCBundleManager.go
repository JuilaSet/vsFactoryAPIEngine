package core

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
)

// 传递RPC参数的对象
type M map[string]interface{}

// 服务注册时请求的表单内容将放入在这个结构体中
type RegisterFrom struct {
	Name string `form:"Name" json:"Name"`
	URL string `form:"URL" json:"URL"`
	Port string `form:"Port" json:"Port"`
	FuncName string `form:"FuncName" json:"FuncName"`
	Describe string `form:"Describe" json:"Describe"`
	FuncNames map[string]bool		// 多个可调用对象
}

// 用于 插件注册 & 插件注销
type RPCBundleManager struct {
	requestMap map[string]IRPCBundleRequest
}

func NewRPCBundleManager() *RPCBundleManager {
	return &RPCBundleManager{requestMap: make(map[string]IRPCBundleRequest, 10)}
}

// 将HttpContext对象升级到Request对象, 并加入注册表中
func (m *RPCBundleManager) Bind(c *gin.Context) (err error) {
	req, err := NewRPCBundleRequest(c)
	m.requestMap[req.form.Name] = req
	// 检测RPC连接
	funcName := req.form.Name + ".Test"
	reply := req.DialRPC(funcName, M{})
	if reply.String() != "OK" {
		err = errors.New("RPC Test failed: " + funcName)
	}
	fmt.Println("RPC test success for " + req.form.Name)
	return
}

// "/" + 函数名
func (m *RPCBundleManager) RPCServe(serviceGroup *gin.RouterGroup, serviceName string, params map[string]interface{}, callBack func(*gin.Context, IRPCBundleRequest, IRPCBundleResponse)) error {
	req, ok := m.requestMap[serviceName]
	if !ok {
		return errors.New("no service named: " + serviceName)
	}
	for suffix := range req.FuncNames() {
		funcName := serviceName + "." + suffix
		fmt.Println("suffix: " + funcName)
		serviceGroup.POST("/" + suffix, func(c *gin.Context) {
			reply := req.DialRPC(funcName, params)
			fmt.Println(reply.String())
			callBack(c, req, reply)
		})
	}

	return nil
}

// urlPath是不带/号的
func (m *RPCBundleManager) DialRPC(serviceName string, urlPath string, params map[string]interface{}) (IRPCBundleResponse, error) {
	req, ok := m.requestMap[serviceName]
	if !ok {
		return nil, errors.New("no service named: " + serviceName)
	}
	funcName := serviceName + "." + urlPath
	reply := req.DialRPC(funcName, params)
	// 获取RPC响应对象
	return reply, nil
}

type IRPCBundleRequest interface{
	FuncNames() map[string]bool
	ContainFunc(funcName string) bool
	DialRPC(funcName string, params map[string]interface{}) IRPCBundleResponse
}

type IRPCBundleResponse interface {
	String() string
}

// RPC服务注册请求对象
type RPCBundleRequest struct {
	form RegisterFrom
}

// RPC服务注册请求对象
type RPCBundleResponse struct {
	responseString string
}

func NewRPCBundleResponse(resp string) *RPCBundleResponse {
	return &RPCBundleResponse{resp}
}

// 获取返回的数据
func (res *RPCBundleResponse) String() string {
	return res.responseString
}

// 构造: RPC服务注册请求对象
func NewRPCBundleRequest(c *gin.Context) (*RPCBundleRequest, error) {
	var form RegisterFrom
	// Bind()默认解析并绑定form格式
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	c.JSON(http.StatusOK, gin.H{"status": "200"})

	funcNames := strings.Split(form.FuncName, ";")
	form.FuncNames = make(map[string]bool)
	for _, v := range funcNames {
		form.FuncNames[v] = true
	}
	return &RPCBundleRequest{form}, nil
}

// 获取所有的方法对象
func (req *RPCBundleRequest) FuncNames() map[string]bool {
	return req.form.FuncNames
}

// 是否包含函数
func (req *RPCBundleRequest) ContainFunc(funcName string)bool {
	_, ok := req.form.FuncNames[funcName]
	return ok
}

// 发送JSON-RPC调用
func (req *RPCBundleRequest) DialRPC(funcName string, params map[string]interface{}) IRPCBundleResponse {
	conn, err := net.Dial("tcp", req.form.URL + req.form.Port)
	if err != nil{
		panic(err)
	}
	// 建立基于json编解码的rpc服务
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply string
	// 调用rpc服务方法
	err = client.Call(funcName, params, &reply)
	if err != nil {
		panic(err)
	}
	return NewRPCBundleResponse(reply)
}